package repository

import (
	"lang/pkg/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type ChapterPost interface {
	CreateChapter(chapter models.Chapter) (int, error)
}

type Repository struct {
	Authorization
	ChapterPost
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		ChapterPost:   NewChapterPostgres(db),
	}
}
