package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"lang/pkg/models"
	"log"
	"strings"
	"time"
)

type ArticlePostgres struct {
	db *sqlx.DB
}

func NewArticlePostgres(db *sqlx.DB) *ArticlePostgres {
	return &ArticlePostgres{db: db}
}

func (r *ArticlePostgres) CreateArticle(article *models.Article) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error with begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("error while rollback transaction: %v", rbErr)
			}
		}
	}()

	var id int

	article.CreatedAt = time.Now()

	createArticleQuery := fmt.Sprintf("INSERT INTO %s (chapter_id, title, level_id, description, image_url, image_alt, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING article_id;", articleTable)

	err = tx.QueryRow(createArticleQuery, article.ChapterID, article.Title, article.LevelId, article.Description, article.ImageUrl, article.ImageAlt, article.CreatedAt).Scan(&id)

	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("error with commit transaction: %w", err)
	}

	return id, nil
}

func (r *ArticlePostgres) GetArticleById(id int) (models.Article, error) {

	var article models.Article

	query := fmt.Sprintf("SELECT chapter_id, article_id, title, level_id, description, image_url, image_alt  FROM %s WHERE article_id = $1", articleTable)
	err := r.db.Get(&article, query, id)

	if err != nil {
		return article, err
	}

	return article, nil
}

func (r *ArticlePostgres) GetAllArticles(id int) ([]models.Article, error) {
	var articles []models.Article

	query := fmt.Sprintf("SELECT article_id, title, description,chapter_id, level_id, image_url, image_alt FROM %s WHERE chapter_id = $1", articleTable)
	err := r.db.Select(&articles, query, id)
	if err != nil {
		return articles, err
	}

	return articles, nil
}

func (r *ArticlePostgres) UpdateArticle(articleId int, article models.UpdateArticle) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	// Проверяем поля на обновление и формируем запрос
	if article.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *article.Title)
		argId++
	}

	if article.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *article.Description)
		argId++
	}

	if article.LevelId != nil {
		setValues = append(setValues, fmt.Sprintf("level_id=$%d", argId))
		args = append(args, *article.LevelId)
		argId++
	}

	if article.ImageUrl != nil {
		setValues = append(setValues, fmt.Sprintf("image_url=$%d", argId))
		args = append(args, *article.ImageUrl)
		argId++
	}

	if article.ImageAlt != nil {
		setValues = append(setValues, fmt.Sprintf("image_alt=$%d", argId))
		args = append(args, *article.ImageAlt)
		argId++
	}

	// Установка времени обновления
	article.UpdateAt = time.Now()
	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, article.UpdateAt)
	argId++

	// Обработка случая, когда нет значений для обновления
	if len(setValues) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE article_id = $%d", articleTable, setQuery, argId)
	args = append(args, articleId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args)

	// Выполнение запроса
	_, err := r.db.Exec(query, args...)
	if err != nil {
		logrus.Errorf("error executing update: %s", err)
		return err
	}

	return nil

}

func (r *ArticlePostgres) Delete(articleId int) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE article_id = $1", articleTable)
	err := r.db.QueryRow(query, articleId)
	if err != nil {
		return err.Err()
	}

	return nil
}
