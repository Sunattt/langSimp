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
	chapter.CreatedAt = time.Now()

	createChapterQuery := fmt.Sprintf("INSERT INTO %s (language_id, title, description, image_url, image_alt, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING chapter_id", chapterTable)

	row := tx.QueryRow(createChapterQuery, chapter.LanguageId, chapter.Title, chapter.Description, chapter.ImageUrl, chapter.ImageAlt, chapter.CreatedAt)
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

	// Проверяем поля на обновление и формируем запрос
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

	// Установка времени обновления
	chp.UpdateAt = time.Now()
	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, chp.UpdateAt)
	argId++

	// Обработка случая, когда нет значений для обновления
	if len(setValues) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Формирование запроса
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE chapter_id = $%d", chapterTable, setQuery, argId)
	args = append(args, chapterId)

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

func (r *ChapterPostgres) GetALL(langId int) ([]models.Chapter, error) {
	var chapters []models.Chapter

	query := fmt.Sprintf(`
        SELECT  ch.chapter_id, ch.language_id, ch.title, ch.description, ch.image_url,  ch.image_alt,
        (SELECT COUNT(*) FROM %s a WHERE a.chapter_id = ch.chapter_id AND a.active = true) AS count_articles
        FROM %s ch WHERE ch.active = true AND ch.language_id = $1
        ORDER BY ch.chapter_id`, // Добавлена сортировка для стабильного порядка
		articleTable,
		chapterTable)

	err := r.db.Select(&chapters, query, langId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}

	if len(chapters) == 0 {
		return []models.Chapter{}, nil // Возвращаем пустой slice вместо nil
	}

	return chapters, nil
}

func (r *ChapterPostgres) GetChapterById(chapterId, langId int) (models.Chapter, error) {
	var chap models.Chapter

	// Основной запрос для получения данных главы
	query := fmt.Sprintf(` SELECT c.chapter_id, c.title, c.description, c.image_url, c.image_alt,
            (SELECT COUNT(*) FROM %s a WHERE a.chapter_id = c.chapter_id AND a.active = true) as count_articles FROM %s c
        WHERE c.chapter_id = $1 AND c.language_id = $2`, articleTable, chapterTable)

	err := r.db.Get(&chap, query, chapterId, langId)
	if err != nil {
		return models.Chapter{}, fmt.Errorf("failed to get chapter: %w", err)
	}

	return chap, nil
}

func (r *ChapterPostgres) Delete(chapterId int) error {
	query := fmt.Sprintf("DELETE FROM %s chp WHERE chp.chapter_id = $1 CASCADE;", chapterTable)
	_, err := r.db.Exec(query, chapterId)
	return err
}
