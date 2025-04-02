package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"lang/pkg/models"
)

type ArticlePostgres struct {
	db *sqlx.DB
}

func NewArticlePostgres(db *sqlx.DB) *ArticlePostgres {
	return &ArticlePostgres{db: db}
}

func (r *ArticlePostgres) createArticle(article *models.Article) (int, error) {
	var id int

	createArticleQuery := fmt.Sprintf("INSERT INTO %s (chapter_id, title, slug, description, image_url, image_alt, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id;", articleTable)

	err := r.db.QueryRow(createArticleQuery, article.ChapterID, article.Title, article.s)
}
