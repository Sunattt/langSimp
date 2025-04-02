package service

import (
	"lang/internal/repository"
	"lang/pkg/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerationToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Verification interface {
	GetUserActive(userId int, username string) (bool, error)
}

type ChapterPost interface {
	Create(chapter *models.Chapter) (int, error)
	GetALL() ([]models.Chapter, error)
	GetChapterById(chapterId int) (models.Chapter, error)
	Update(chapterId int, input models.UpdateChapter) error
}

type Service struct {
	Authorization
	ChapterPost
	Verification
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, repo.Verification),
		ChapterPost:   NewChapterService(repo.ChapterPost),
		Verification:  NewVerService(repo.Verification),
	}
}
