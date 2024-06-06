package main

import (
	"log"
	"net/http"
	"vectorsearch-go/domain"
	"vectorsearch-go/infra"

	"github.com/labstack/echo/v4"
)

func main() {
	es_article_repository, err := infra.NewESArticleRepository()
	if err != nil {
		log.Fatal(err)
		return
	}
	ai, err := infra.NewGenAi()
	if err != nil {
		log.Fatal(err)
		return
	}

	service := NewService(es_article_repository, ai)
	e := echo.New()

	e.GET("/articles", service.GetArticles)
	e.GET("/article/:id", service.GetArticle)
	e.POST("/article/:id", service.CreateArticle)
	e.GET("/search", service.SearchArticle)
	e.GET("/semantics", service.SemanticsSearch)

	e.Logger.Fatal(e.Start(":8080"))

}

type Service struct {
	ArticleRepository domain.ArticleRepository
	Ai                GenAi
}

func NewService(ar domain.ArticleRepository, ai GenAi) *Service {
	return &Service{
		ArticleRepository: ar,
		Ai:                ai,
	}
}

func (s *Service) GetArticles(c echo.Context) error {
	articles, err := s.ArticleRepository.List()
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, articles)
}

func (s *Service) CreateArticle(c echo.Context) error {
	a := new(struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	})
	if err := c.Bind(a); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	article, err := domain.NewArticle(c.Param("id"), a.Title, a.Content)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	embedding, err := s.Ai.Embedding(article.Content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := s.ArticleRepository.Store(article.Title, article.Content, embedding); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(200, article)
}

func (s *Service) GetArticle(c echo.Context) error {
	article, err := s.ArticleRepository.Find(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(200, article)
}

func (s *Service) SearchArticle(c echo.Context) error {
	articles, err := s.ArticleRepository.SearchContent(c.QueryParam("keyword"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(200, articles)
}
func (s *Service) SemanticsSearch(c echo.Context) error {
	text := c.QueryParam("text")
	if text == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "text is required")
	}
	embedding, err := s.Ai.Embedding(text)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	articles, err := s.ArticleRepository.SemanticsSearch(embedding)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(200, articles)
}

type GenAi interface {
	Embedding(text string) ([]float32, error)
}
