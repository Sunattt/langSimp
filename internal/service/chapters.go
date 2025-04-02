package service

import (
	"lang/internal/repository"
	"lang/pkg/models"
)

type ChapterService struct {
	repo repository.ChapterPost
}

func NewChapterService(repo repository.ChapterPost) *ChapterService {
	return &ChapterService{repo: repo}
}

func (s *ChapterService) Create(chap *models.Chapter) (int, error) {
	return s.repo.Create(chap)
}

func (s *ChapterService) GetALL() ([]models.Chapter, error) {
	return s.repo.GetALL()
}

func (s *ChapterService) GetChapterById(chapterId int) (models.Chapter, error) {
	return s.repo.GetChapterById(chapterId)
}

func (s *ChapterService) Update(chapterId int, chap models.UpdateChapter) error {
	return s.repo.Update(chapterId, chap)
}

func (s *ChapterService) Delete(chapterId int) error {
	return s.repo.Delete(chapterId)
}
