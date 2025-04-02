package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"lang/pkg/models"
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
	var id int

	chapter.CreatedAt = time.Now()

	count, err := countArticle(r)
	if err != nil {
		return 0, err
	}

	createChapterQuery := fmt.Sprintf("INSERT INTO %s (language_id, count_article ,title, description, image_url, image_alt, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", chapterTable)

	row := r.db.QueryRow(createChapterQuery, 1, count, chapter.Title, chapter.Description, chapter.ImageUrl, chapter.ImageAlt, chapter.CreatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func countArticle(r *ChapterPostgres) (int, error) {
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s ar INNER JOIN %s chp on ar.article_id = chp.chapter_id WHERE ar.active = true", articleTable, chapterTable)

	if err := r.db.QueryRow(query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
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

	if chp.Description == nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *chp.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s chp SET %s FROM %s WHERE chp.id = $1",
		chapterTable, setQuery, chapterTable)

	args = append(args, chapterId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args)

	_, err := r.db.Exec(query, args...)
	return err

}

func (r *ChapterPostgres) GetALL() ([]models.Chapter, error) {
	var chapters []models.Chapter

	query := fmt.Sprintf("SELECT chp.chapter_id, chp.title, chp.description, chp.image, chp.count_article FROM %s chp WHERE active = true", chapterTable)
	err := r.db.Select(&chapters, query)

	return chapters, err
}

func (r *ChapterPostgres) GetChapterById(chapterId int) (models.Chapter, error) {
	var chap models.Chapter

	query := fmt.Sprintf(" SELECT chp.chapter_id, chp.title, chp.description, chp.image, chp.count_article FROM %s chp WHERE chp.chapter_id = $1", chapterTable)

	err := r.db.Get(&chap, query, chapterId)

	return chap, err
}

func (r *ChapterPostgres) Delete(chapterId int) error {
	query := fmt.Sprintf("DELETE FROM %s chp WHERE chp.chapter_id = $1")
	_, err := r.db.Exec(query, chapterId)
	return err
}
