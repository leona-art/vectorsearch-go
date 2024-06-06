package domain

import "fmt"

type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewArticle(id, title, content string) (*Article, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}
	return &Article{
		Id:      id,
		Title:   title,
		Content: content,
	}, nil
}

func (a *Article) GetID() string {
	return a.Id
}

func (a *Article) GetTitle() string {
	return a.Title
}

func (a *Article) GetContent() string {
	return a.Content
}

type ArticleWithScore struct {
	Article Article
	Score   float64
}

func NewArticleWithScore(article Article, score float64) *ArticleWithScore {
	return &ArticleWithScore{
		Article: article,
		Score:   score,
	}
}
