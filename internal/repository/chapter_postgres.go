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

type ChapterPostgres struct {
	db *sqlx.DB
}

func NewChapterPostgres(db *sqlx.DB) *ChapterPostgres {
	return &ChapterPostgres{db: db}
}

const (
	articleTable = "articles"
	chapterTable = "chapters"
)

func (r *ChapterPostgres) Create(chapter *models.Chapter) (int, error) {
	tx, err := r.db.Begin() // Исправлено с txt на tx
	if err != nil {
		return 0, fmt.Errorf("error with begin transaction: %w", err)
	}

	defer func() {
		if err != nil { // Проверяем, была ли ошибка, чтобы выполнить Rollback
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("error while rollback transaction: %v", rbErr)
			}
		}
	}()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s ar INNER JOIN %s chp ON ar.chapter_id = chp.chapter_id WHERE ar.active = true", articleTable, chapterTable)

	if err := tx.QueryRow(query).Scan(&count); err != nil {
		log.Println("while counting articles in chapter:", err.Error())
		return 0, err
	}

	var id int
	chapter.CreatedAt = time.Now()

	createChapterQuery := fmt.Sprintf("INSERT INTO %s (language_id, count_articles, title, description, image_url, image_alt, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING chapter_id", chapterTable)

	row := tx.QueryRow(createChapterQuery, chapter.LanguageId, count, chapter.Title, chapter.Description, chapter.ImageUrl, chapter.ImageAlt, chapter.CreatedAt)
	if err := row.Scan(&id); err != nil {
		log.Println("while saving chapter:", err.Error())
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("while committing chapter:", err.Error())
		return 0, err
	}

	return id, nil
}

func (r *ChapterPostgres) Update(chapterId int, chp models.UpdateChapter) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if chp.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *chp.Title)
		argId++
	}

	if chp.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *chp.Description)
		argId++
	}

	if chp.ImageUrl != nil {
		setValues = append(setValues, fmt.Sprintf("image_url=$%d", argId))
		args = append(args, *chp.ImageUrl)
		argId++
	}

	if chp.ImageAlt != nil {
		setValues = append(setValues, fmt.Sprintf("image_alt=$%d", argId))
		args = append(args, *chp.ImageAlt)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	chp.UpdateAt = time.Now()

	query := fmt.Sprintf("UPDATE %s chp SET %s FROM %s WHERE chp.id = $1",
		chapterTable, setQuery, chapterTable)

	args = append(args, chapterId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args)

	_, err := r.db.Exec(query, args...)
	return err

}

func (r *ChapterPostgres) GetALL(langId int) ([]models.Chapter, error) {
	var chapters []models.Chapter

	query := fmt.Sprintf("SELECT chp.chapter_id, chp.title, chp.description, chp.image_url, chp.image_alt chp.count_articles FROM %s chp WHERE active = true AND language_id = $2", chapterTable)
	err := r.db.Select(&chapters, query, langId)

	return chapters, err
}

func (r *ChapterPostgres) GetChapterById(chapterId int) (models.Chapter, error) {
	var chap models.Chapter

	query := fmt.Sprintf(" SELECT chp.chapter_id, chp.title, chp.description, chp.image_url, chp.image_alt, chp.count_articles FROM %s chp WHERE chp.chapter_id = $1", chapterTable)

	err := r.db.Get(&chap, query, chapterId)

	return chap, err
}

func (r *ChapterPostgres) Delete(chapterId int) error {
	query := fmt.Sprintf("DELETE FROM %s chp WHERE chp.chapter_id = $1", chapterTable)
	_, err := r.db.Exec(query, chapterId)
	return err
}
