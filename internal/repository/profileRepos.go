package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"lang/pkg/models"
	"time"
)

type ProfileRepos struct {
	db *sqlx.DB
}

func NewProfileRepos(db *sqlx.DB) *ProfileRepos {
	return &ProfileRepos{db: db}
}

func (r *ProfileRepos) SaveChapter(chapterId, userId int) error {

	savedAt := time.Now()

	query := fmt.Sprintf(`INSERT INTO user_saved_chapters (user_id, chapter_id, saved_at)
        VALUES ($1, $2, 3$)`)

	err := r.db.QueryRow(query, chapterId, userId, savedAt)
	if err != nil {
		return fmt.Errorf("Error Save Chapter %v ", err)
	}
	return nil
}

func (r *ProfileRepos) RemoveSavedChapter(chapterId, userId int) error {

	query := fmt.Sprintf(`DELETE FROM user_saved_chapters WHERE user_id = $1 AND chapter_id = $2`)

	err := r.db.QueryRow(query, userId, chapterId)
	if err != nil {
		return fmt.Errorf("Error Remove Chapter %v ", err)
	}

	return nil
}

func (r *ProfileRepos) GetSavedChapters(userId int) ([]models.Chapter, error) {

	query := fmt.Sprintf(`SELECT ch.chapter_id,ch.language_id, (SELECT COUNT(*) FROM articles a WHERE a.chapter_id = ch.chapter_id AND a.active = true) AS ch.count_articles,ch.title,ch.description, ch.image_url, 
      	 ch.image_alt, ch.created_at, uc.saved_at FROM chapters ch JOIN user_saved_chapters uc ON uc.chapter_id = ch.chapter_id 
			WHERE uc.user_id = $1 AND ch.active = true ORDER BY uc.saved_at DESC`)

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("Error Get Saved Chapter %v ", err)
	}

	defer rows.Close()

	var savedChapters []models.Chapter
	for rows.Next() {
		var sch models.Chapter
		err := rows.Scan(
			&sch.Id,
			&sch.LanguageId,
			&sch.CountArticle,
			&sch.Title,
			&sch.Description,
			&sch.ImageUrl,
			&sch.ImageAlt,
			&sch.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("Error Get Saved Chapter %v ", err)
		}

		savedChapters = append(savedChapters, sch)
	}
	return savedChapters, nil
}

func (r *ProfileRepos) SaveArticle(articleId, userId int) error {
	savedAt := time.Now()

	query := fmt.Sprintf(`INSERT INTO user_saved_articles (user_id, article_id, saved_at)
        VALUES ($1, $2, 3$)`)

	err := r.db.QueryRow(query, articleTable, userId, savedAt)
	if err != nil {
		return fmt.Errorf("Error Save Chapter %v ", err)
	}
	return nil
}

func (r *ProfileRepos) RemoveSavedArticle(chapterId, userId int) error {

	query := fmt.Sprintf("DELETE FROM user_saved_articles where user_id = $1 AND chapter_id = $2")

	err := r.db.QueryRow(query, userId, chapterId)
	if err != nil {
		return fmt.Errorf("Error Remove Article %v ", err)
	}
	return nil
}

func (r *ProfileRepos) GetSavedArticles(userId int) ([]models.Article, error) {

	query := `SELECT a.article_id, a.chapter_id, a.title, a.level_id, a.description, a.image_url, 
       			a.image_alt, a.created_at, uc.saved_at FROM articles a JOIN user_saves_articles ua ON ua.article_id = a.article_id 
       			WHERE a.user_id = $1 AND a.active = true  ORDER BY ua.saved_at DESC`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("Error Get Saved Article %v ", err)
	}

	defer rows.Close()

	var savedArticles []models.Article
	for rows.Next() {
		var art models.Article
		err := rows.Scan(
			&art.Id,
			&art.ChapterID,
			&art.Title,
			&art.LevelId,
			&art.Description,
			&art.ImageUrl,
			&art.ImageAlt,
			&art.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("Error Get Saved Article %v ", err)
		}

		savedArticles = append(savedArticles, art)
	}
	return savedArticles, nil
}

func (r *ProfileRepos) SaveWord(wordId, userId int) error {
	return nil
}
