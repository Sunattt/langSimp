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
	IsAdmin(userId int) (bool, error)
	IsModerator(userId int) (bool, error)
	IsEmailFree(email string) (bool, error)
}

type ChapterPost interface {
	Create(chapter *models.Chapter) (int, error)
	GetALL(landId int) ([]models.Chapter, error)
	GetChapterById(chapterId int) (models.Chapter, error)
	Update(chapterId int, input models.UpdateChapter) error
	Delete(chapterId int) error
}

type ActiclePost interface {
	Create(article *models.Article) (int, error)
	GetALL(landId int) ([]models.Article, error)
	GetChapterById(article int) (models.Article, error)
	Update(chapterId int, chp models.UpdateChapter) error
	Delete(chapterId int) error
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
