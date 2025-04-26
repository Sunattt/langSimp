package service

import (
	"lang/internal/repository"
	"lang/pkg/models"
	"net/http"
	"strings"
)

type ArticleService struct {
	repo repository.ArticlePost
	ver  repository.Verification
}

func NewArticleService(repo repository.ArticlePost) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) CreateArticle(article *models.Article) (int, models.ErrorResponse) {

	if strings.TrimSpace(article.Title) == "" || strings.TrimSpace(article.Description) == "" {
		return 0, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Article content cannot be empty",
		}
	}

	exist, err := s.ver.IsValidChapterId(article.ChapterID)
	if !exist {
		return 0, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid chapter id",
			Error:   err,
		}
	}

	articleId, err := s.repo.CreateArticle(article)
	if err != nil {
		return 0, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "error while creating article",
			Error:   err,
		}
	}

	return articleId, models.ErrorResponse{}
}

func (s *ArticleService) GetAllArticles(langId int) ([]models.Article, error) {
	return s.repo.GetAllArticles(langId)
}

func (s *ArticleService) GetArticleById(articleId int) (models.Article, error) {
	return s.repo.GetArticleById(articleId)
}

func (s *ArticleService) Update(articleId int, chp models.UpdateArticle) error {
	return s.repo.UpdateArticle(articleId, chp)
}

func (s *ArticleService) Delete(articleId int) error {
	return s.repo.Delete(articleId)
}
