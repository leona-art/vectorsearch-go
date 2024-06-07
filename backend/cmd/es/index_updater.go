package main

import (
	"context"
	"vectorsearch-go/infra"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const (
	ES_ADDRESS = "http://localhost:9200"
)

func main() {

	config := elasticsearch.Config{
		Addresses: []string{
			ES_ADDRESS,
		},
	}
	es, err := elasticsearch.NewTypedClient(config)
	if err != nil {
		panic(err)

	}
	is_exists, err := es.Indices.Exists(infra.ES_INDEX).Do(context.Background())

	if err != nil {
		panic(err)
	}
	if is_exists {
		_, err := es.Indices.Delete(infra.ES_INDEX).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	_, create_err := es.Indices.Create(infra.ES_INDEX).
		Request(&create.Request{
			// Settings: &types.IndexSettings{
			// 	Analysis: &types.IndexSettingsAnalysis{
			// 		Tokenizer: map[string]types.Tokenizer{
			// 			"kromoji_tokenizer": types.KuromojiTokenizer{
			// 				Type: "kuromoji_tokenizer",
			// 				Mode: kuromojitokenizationmode.Normal,
			// 			},
			// 		},
			// 		Analyzer: map[string]types.Analyzer{
			// 			"kromoji_analyzer": types.KuromojiAnalyzer{
			// 				Type: "custom",
			// 				Mode: kuromojitokenizationmode.Normal,
			// 			},
			// 		},
			// 	},
			// },
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
	if create_err != nil {
		panic(create_err)
	}
}
