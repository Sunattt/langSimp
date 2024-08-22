package repository

import (
	"github.com/jmoiron/sqlx"
	"lang/pkg/models"
)

type ChapterPostgres struct {
	db *sqlx.DB
}

func NewChapterPostgres(db *sqlx.DB) *ChapterPostgres {
	return &ChapterPostgres{db: db}
}
func (r *ChapterPostgres) CreateChapter(chapter models.Chapter) (int, error) {
	//TODO implement me
	panic("implement me")
}
