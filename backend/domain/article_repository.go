package domain

type ArticleRepository interface {
	Store(title string, content string, vector []float32) error
	Find(id string) (*Article, error)
	List() ([]*Article, error)
	SearchContent(keyword string) ([]*Article, error)
	SemanticsSearch(vector []float32) ([]*ArticleWithScore, error)
}
