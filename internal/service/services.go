package service

import (
	"lang/internal/repository"
	"lang/pkg/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerationToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	UserDataValidation(user models.User) error
}

type ChapterPost interface {
	CreateChapter(chapter models.Chapter) (int, error)
}

type Service struct {
	Authorization
	ChapterPost
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
