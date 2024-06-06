package infra

import (
	"context"
	"errors"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GoogleAi struct {
}

const (
	MODEL_CODE     = "models/text-embedding-004"
	EMBEDDING_DIMS = 768
)

func NewGenAi() (*GoogleAi, error) {

	return &GoogleAi{}, nil
}

type Embeddings [EMBEDDING_DIMS]float32

func NewEmbeddings(e []float32) (*Embeddings, error) {
	if len(e) != EMBEDDING_DIMS {
		return nil, errors.New("vector length must be 768")
	}
	embed := Embeddings(e)
	return &embed, nil
}
func (e *Embeddings) ToSlice() []float32 {
	return e[:]
}

func (g *GoogleAi) Embedding(text string) ([]float32, error) {
	// Set up a context and a client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GOOGLE_AI_API_KEY")))
	if err != nil {
		return nil, err
	}

	defer client.Close()

	// Embed the text
	em := client.EmbeddingModel(MODEL_CODE)
	res, err := em.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, err
	}
	embedding, err := NewEmbeddings(res.Embedding.Values)
	if err != nil {
		return nil, err
	}
	return embedding.ToSlice(), nil
}
