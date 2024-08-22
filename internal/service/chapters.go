package service

import (
	"lang/internal/repository"
	"lang/pkg/models"
)

type ChapterService struct {
	repo repository.ChapterPostgres
}

func NewChapterService(repo repository.ChapterPostgres) *ChapterService {
	return &ChapterService{repo: repo}
}

func (s *ChapterService) CreateChapter(chap models.Chapter) (int, error) {
	return s.repo.CreateChapter(chap)
}
