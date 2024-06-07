package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
	"vectorsearch-go/domain"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const (
	ES_INDEX    = "articles"
	VECTOR_DIMS = 768
)

type esArticleRepository struct {
	client *elasticsearch.TypedClient
}

func retry[T any](attempts int, sleep time.Duration, f func() (*T, error)) (*T, error) {
	var err error
	var res *T
	for i := 0; i < attempts; i++ {
		if res, err = f(); err == nil {
			return res, nil
		}
		fmt.Printf("retrying after error: %v", err)
		time.Sleep(sleep)
	}
	return nil, fmt.Errorf("failed after %d attempts, last error: %v", attempts, err)
}

func NewESArticleRepository() (*esArticleRepository, error) {
	esUrl := os.Getenv("ELASTICSEARCH_URL")
	if esUrl == "" {
		return nil, errors.New("ELASTICSEARCH_URL is required")
	}
	config := elasticsearch.Config{
		Addresses: []string{
			esUrl,
		},
	}

	es, err := retry(3, 15*time.Second, func() (*elasticsearch.TypedClient, error) {
		return elasticsearch.NewTypedClient(config)
	})
	fmt.Println("Connected to Elastic Search")
	if err != nil {
		return nil, err
	}

	exists_index, err := retry(3, 30*time.Second, func() (*bool, error) {
		res, err := es.Indices.Exists(ES_INDEX).Do(context.Background())
		if err != nil {
			return nil, err
		}
		return &res, nil
	})

	if err != nil {
		return nil, err
	}
	if !*exists_index {
		_, err := es.Indices.
			Create(ES_INDEX).
			Request(&create.Request{
				Mappings: &types.TypeMapping{
					Properties: map[string]types.Property{
						"title":   types.NewTextProperty(),
						"content": types.NewKeywordProperty(),
						"embed": types.DenseVectorProperty{
							Dims: &[]int{768}[0],
							Type: "dense_vector",
						},
					},
				},
			}).
			Do(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return &esArticleRepository{
		client: es,
	}, nil
}

type ESArticle struct {
	Title   string               `json:"title"`
	Content string               `json:"content"`
	Embed   [VECTOR_DIMS]float32 `json:"embed"`
}

func NewESArticle(title string, content string, vector []float32) (*ESArticle, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if len(vector) != VECTOR_DIMS {
		return nil, errors.New("vector length must be 768")
	}
	return &ESArticle{
		Title:   title,
		Content: content,
		Embed:   [VECTOR_DIMS]float32(vector),
	}, nil
}

func (r *esArticleRepository) Store(title string, content string, vector []float32) error {
	article, err := NewESArticle(title, content, vector)
	if err != nil {
		return err
	}

	if _, err := r.client.
		Index(ES_INDEX).
		Document(article).
		Do(context.Background()); err != nil {
		return err
	}
	return nil
}

type FindArticleResponse struct {
	Index       string    `json:"_index"`
	Id          string    `json:"_id"`
	Version     int       `json:"_version"`
	SeqNo       int       `json:"_seq_no"`
	PrimaryTerm int       `json:"_primary_term"`
	Found       bool      `json:"found"`
	Source      ESArticle `json:"_source"`
}

func (r *esArticleRepository) Find(id string) (*domain.Article, error) {
	res, err := r.client.Get(ES_INDEX, id).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	var article ESArticle
	if err := json.Unmarshal(res.Source_, &article); err != nil {
		return nil, err
	}
	return domain.NewArticle(res.Id_, article.Title, article.Content)
}

type ArticleSearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float32 `json:"max_score"`
		Hits     []struct {
			Index  string    `json:"_index"`
			Type   string    `json:"_type"`
			Id     string    `json:"_id"`
			Score  float32   `json:"_score"`
			Source ESArticle `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (r *esArticleRepository) List() ([]*domain.Article, error) {
	res, err := r.client.Search().
		Index(ES_INDEX).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	articles := make([]*domain.Article, 0)
	for _, hit := range res.Hits.Hits {
		var es_article ESArticle
		if err := json.Unmarshal(hit.Source_, &es_article); err != nil {
			return nil, err
		}
		article, err := domain.NewArticle(hit.Id_, es_article.Title, es_article.Content)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (r *esArticleRepository) SearchContent(keyword string) ([]*domain.Article, error) {
	res, err := r.client.Search().
		Index(ES_INDEX).
		Request(&search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"content": {
						Query: keyword,
					},
				},
			},
		}).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	articles := make([]*domain.Article, 0)
	for _, hit := range res.Hits.Hits {
		var es_article ESArticle
		if err := json.Unmarshal(hit.Source_, &es_article); err != nil {
			return nil, err
		}
		article, err := domain.NewArticle(hit.Id_, es_article.Title, es_article.Content)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (r *esArticleRepository) SemanticsSearch(vector []float32) ([]*domain.ArticleWithScore, error) {
	res, err := r.client.Search().
		Index(ES_INDEX).
		Knn(types.KnnQuery{
			Field:         "embed",
			QueryVector:   vector,
			K:             10,
			NumCandidates: 100,
		}).
		Fields(types.FieldAndFormat{
			Field: "title",
		}).
		Fields(types.FieldAndFormat{
			Field: "content",
		}).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	articles := make([]*domain.ArticleWithScore, 0)
	for _, hit := range res.Hits.Hits {
		var es_article ESArticle

		if err := json.Unmarshal(hit.Source_, &es_article); err != nil {
			return nil, err
		}
		article, err := domain.NewArticle(hit.Id_, es_article.Title, es_article.Content)
		if err != nil {
			return nil, err
		}
		articles = append(articles, domain.NewArticleWithScore(*article, float64(hit.Score_)))
	}
	return articles, nil
}
