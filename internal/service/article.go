package service

import (
	"lang/internal/repository"
	"lang/pkg/models"
)

type ArticleService struct {
	repo repository.ChapterPost
}

func NewArticleService(repo repository.ChapterPost) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) createArticle(article *models.Article) (int, error) {
	return s.repo.Create(article)
}
