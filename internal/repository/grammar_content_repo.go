package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"lang/pkg/models"
	"log"
	"strings"
	"time"
)

const (
	grammarConTable     = "grammar_contents"
	grammarExrTable     = "grammar_exercises"
	grammarCommentTable = "grammar_comments"
	commentLikes        = "comment_likes"
)

type ContentPostgres struct {
	db *sqlx.DB
}

func NewContentPostgres(db *sqlx.DB) *ContentPostgres {
	return &ContentPostgres{db: db}
}

func (r *ContentPostgres) CreateContent(input *models.GrammarContent) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error with begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("error while rollback transaction: %v ", rbErr)
			}
		}
	}()

	jsonExample, err := json.Marshal(input.Example)
	if err != nil {
		return 0, fmt.Errorf("error with marshalling input: %w", err)
	}

	var id int
	input.CreatedAt = time.Now()
	query := fmt.Sprintf(`INSERT INTO %s (article_id, level_id, explanation, structure, examples, tips, picture, created_at)
	VALUES 	($1, $2, $3, $4,$5,$6,$7, $8) RETURNING id`, grammarConTable)

	err = tx.QueryRow(query, input.ArticleId, input.LevelId, input.Explanation, input.Structure, jsonExample, input.Tips,
		input.Picture, input.CreatedAt).Scan(&id)
	if err != nil {
		tx.Rollback()
		log.Println("error while creating content:", err)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Println("error while commit transaction:", err)
		return 0, err
	}
	return id, nil
}

func (r *ContentPostgres) GetCourseById(contentId, levelId int) (models.GrammarContent, error) {
	query := fmt.Sprintf(`
    SELECT id, article_id, level_id, explanation, structure, examples, tips, picture, created_at 
    FROM %s 
    WHERE article_id = $1 AND level_id = $2`,
		grammarConTable)

	var rawExample []byte // Временная переменная для хранения JSONB данных
	input := models.GrammarContent{}

	// Сканируем данные, сохраняя examples как []byte
	err := r.db.QueryRowx(query, contentId, levelId).Scan(
		&input.Id,
		&input.ArticleId,
		&input.LevelId,
		&input.Explanation,
		&input.Structure,
		&rawExample, // Получаем сырые JSON-данные
		&input.Tips,
		&input.Picture,
		&input.CreatedAt,
	)

	if err != nil {
		return models.GrammarContent{}, fmt.Errorf("error while getting grammar content: %w", err)
	}

	// Десериализуем JSON в структуру
	if len(rawExample) > 0 {
		if err := json.Unmarshal(rawExample, &input.Example); err != nil {
			return models.GrammarContent{}, fmt.Errorf("error unmarshaling examples: %w", err)
		}
	}

	return input, nil
}

func (r *ContentPostgres) UpdateContent(contentId, levelId int, input models.UpdateContentInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Explanation != nil {
		setValues = append(setValues, fmt.Sprintf("explanation=$%d", argId))
		args = append(args, *input.Explanation)
		argId++
	}

	if input.Structure != nil {
		setValues = append(setValues, fmt.Sprintf("structure=$%d", argId))
		args = append(args, *input.Structure)
		argId++
	}

	// Особенная обработка для JSON поля examples
	if input.Example != nil {
		jsonExample, err := json.Marshal(input.Example)
		if err != nil {
			return fmt.Errorf("failed to marshal examples: %w", err)
		}
		setValues = append(setValues, fmt.Sprintf("examples=$%d::jsonb", argId))
		args = append(args, jsonExample)
		argId++
	}

	if input.Tips != nil {
		setValues = append(setValues, fmt.Sprintf("tips=$%d", argId))
		args = append(args, *input.Tips)
		argId++
	}

	if input.Picture != nil {
		setValues = append(setValues, fmt.Sprintf("picture=$%d", argId))
		args = append(args, *input.Picture)
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, time.Now())
	argId++

	if len(setValues) == 0 {
		return errors.New("no fields to update")
	}

	query := fmt.Sprintf(`
        UPDATE %s 
        SET %s 
        WHERE id = $%d AND level_id = $%d`,
		grammarConTable,
		strings.Join(setValues, ", "),
		argId, argId+1)

	// Добавляем ID в конец аргументов
	args = append(args, contentId, levelId)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update grammar content: %w", err)
	}

	return nil
}

func (r *ContentPostgres) DeleteContent(contentId, levelId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND level_id = $2", grammarConTable)
	err := r.db.QueryRow(query, contentId, levelId)
	return err.Err()
}

func (r *ContentPostgres) CreateExercise(exr *models.GrammarContentExercises) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback; err != nil {
				log.Printf("error while rollback transaction: %v ", rbErr)
			}
		}
	}()

	var id int

	query := fmt.Sprintf(`INSERT INTO %s (question, option, correct_answer, explanation, difficulty, help)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`, grammarConTable)

	err = tx.QueryRow(query, exr.Question, exr.QuestionType, exr.Option, exr.CorrectAnswer, exr.Explanation,
		exr.Difficulty, exr.Help).Scan(&id)
	if err != nil {
		tx.Rollback()
		log.Println("error while creating exercise:", err)
		return 0, fmt.Errorf("error while creating exercise: %w", err)
	}

	if err := tx.Commit(); err != nil {
		log.Println("error while commit transaction:", err)
		return 0, err
	}

	return id, nil

}

func (r *ContentPostgres) GetExerciseById(contentId int) (models.GrammarContentExercises, error) {
	var exr models.GrammarContentExercises

	query := fmt.Sprintf(`SELECT id, question, question_type, option, correct_answer, explanation, difficulty, help 
								FROM %s WHERE grammar_content_id = $1`, grammarConTable)

	err := r.db.QueryRowx(query, contentId).StructScan(&exr)
	if err != nil {
		return models.GrammarContentExercises{}, fmt.Errorf("error while getting exercise: %w", err)
	}

	return exr, nil
}

func (r *ContentPostgres) UpdateExercise(contentId int, input models.UpdateGrammarExercise) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Question != nil {
		setValues = append(setValues, fmt.Sprintf("question=$%d", argId))
		args = append(args, *input.Question)
		argId++
	}

	if input.QuestionType != nil {
		setValues = append(setValues, fmt.Sprintf("question_type=$%d", argId))
		args = append(args, *input.QuestionType)
		argId++
	}

	// Особенная обработка для JSON поля examples
	if input.Option != nil {
		jsonExample, err := json.Marshal(input.Option)
		if err != nil {
			return fmt.Errorf("failed to marshal examples: %w", err)
		}
		setValues = append(setValues, fmt.Sprintf("options=$%d::jsonb", argId))
		args = append(args, jsonExample)
		argId++
	}

	if input.CorrectAnswer != nil {
		setValues = append(setValues, fmt.Sprintf("correct_answer=$%d", argId))
		args = append(args, *input.CorrectAnswer)
		argId++
	}

	if input.Explanation != nil {
		setValues = append(setValues, fmt.Sprintf("explanation=$%d", argId))
		args = append(args, *input.Explanation)
		argId++
	}

	if input.Difficulty != nil {
		setValues = append(setValues, fmt.Sprintf("difficulty=$%d", argId))
		args = append(args, *input.Difficulty)
		argId++
	}

	if input.Help != nil {
		setValues = append(setValues, fmt.Sprintf("help=$%d", argId))
		args = append(args, *input.Help)
		argId++
	}

	if len(setValues) == 0 {
		return errors.New("no fields to update")
	}

	query := fmt.Sprintf(`
        UPDATE %s 
        SET %s 
        WHERE grammar_content_id = $%d`,
		grammarExrTable,
		strings.Join(setValues, ", "),
		argId)

	// Добавляем ID в конец аргументов
	args = append(args, contentId)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update grammar exercise content: %w", err)
	}

	return nil
}

func (r *ContentPostgres) DeleteExercise(contentId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE grammar_content_id = $1`, grammarConTable)
	err := r.db.QueryRow(query, contentId)
	if err != nil {
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	return nil
}

func (r *ContentPostgres) CreateComment(userId, contentId int, input models.GrammarComment) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error while starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Println("error while rolling back transaction:", rbErr)
			}
		}
	}()

	var id int

	input.CreatedAt = time.Now()

	query := fmt.Sprintf(`INSERT INTO %s (user_id, grammar_content_id, comment, rating, create_at)	
			VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`, grammarConTable)

	err = tx.QueryRow(query, userId, contentId, input.Comment, input.Rating, input.LikesCount, input.CreatedAt).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error while creating comment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("error while commit transaction: %w", err)
	}

	return input.Id, nil
}

func (r *ContentPostgres) GetAllComments(contentId int) ([]models.GrammarComment, error) {
	var comments []models.GrammarComment

	query := fmt.Sprintf(`SELECT id, user_id, grammar_content_id, comment, rating, likes_count, created_at FROM %s WHERE grammar_content_id = $1`, grammarCommentTable)

	err := r.db.Select(&comments, query, contentId)
	if err != nil {
		return nil, fmt.Errorf("error while getting comments: %w", err)
	}

	return comments, nil
}

func (r *ContentPostgres) UpdateComment(commentId int, input models.UpdateGrammarComment) error {

	query := fmt.Sprintf(`UPDATE %s SET comment = $1, updated_at = NOW() WHERE id = $2`, grammarCommentTable)

	err := r.db.QueryRow(query, commentId, input.Comment)
	if err != nil {
		return fmt.Errorf("error while updating comment: %w", err)
	}
	
	return nil

}

func (r *ContentPostgres) DeleteComment(commentId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE grammar_content_id = $1`, grammarCommentTable)
	err := r.db.QueryRow(query, commentId)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

func (r *ContentPostgres) AddLike(userId, commentId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO grammar_likes (user_id, comment_id) VALUES ($1, $2)`
	_, err = tx.Exec(query, userId, commentId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `UPDATE grammar_comments SET likes_count = likes_count + 1 WHERE id = $1`
	_, err = tx.Exec(query, commentId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ContentPostgres) RemoveLike(userId, commentId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `DELETE FROM comment_likes WHERE user_id = $1 AND comment_id = $2`
	_, err = tx.Exec(query, userId, commentId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `UPDATE grammar_comments SET likes_count = likes_count - 1 WHERE comment_id = $1`
	_, err = tx.Exec(query, commentId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ContentPostgres) CheckLikeExists(userId, commentId int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM grammar_likes WHERE user_id = $1 AND comment_id = $2)`
	err := r.db.Get(&exists, query, userId, commentId)
	return exists, err
}
