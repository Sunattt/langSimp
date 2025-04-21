package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	languageTable = "languages"
	levelTable    = "grammar_levels"
)

type VerPostgres struct {
	db *sqlx.DB
}

func NewVerPostgres(db *sqlx.DB) *VerPostgres {
	return &VerPostgres{db: db}
}

func (r *VerPostgres) IsEmailFree(email string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s u WHERE u.email = $1", userTable)

	var countEmail int
	err := r.db.QueryRow(query, email).Scan(&countEmail)
	if err != nil {
		return false, err
	}

	if countEmail != 0 {
		return false, errors.New("email isn't free")
	}

	return true, nil
}

func (r *VerPostgres) GetUserActive(userId int, username string) (bool, error) {
	query := fmt.Sprintf("SELECT active FROM %s u WHERE u.username = $1 AND u.id = $2 ", userTable)

	var active bool
	err := r.db.QueryRow(query, username, userId).Scan(&active)
	if err != nil {
		return false, err // Возвращаем ошибку
	}

	return active, nil // Возвращаем статус активности пользователя
}

func (r *VerPostgres) IsAdmin(userId int) (bool, error) {
	query := fmt.Sprintf("SELECT role_id FROM %s u WHERE u.id = $1", userTable)

	var roleId int
	err := r.db.QueryRow(query, userId).Scan(&roleId)

	if err != nil {
		return false, err
	}

	if roleId == 3 {
		return true, nil
	}

	return false, nil
}

func (r *VerPostgres) IsModerator(userId int) (bool, error) {
	query := fmt.Sprintf("SELECT role FROM %s u WHERE u.id = $1", userTable)

	var roleId int
	err := r.db.QueryRow(query, userId).Scan(&roleId)
	if err != nil {
		return false, err
	}

	if roleId != 2 {
		return false, nil
	}

	return true, nil
}

func (r *VerPostgres) ValidLangCode(langCode string) (int, error) {
	var langId int

	query := fmt.Sprintf("SELECT language_id FROM %s WHERE code = $1", languageTable)

	_ = r.db.QueryRow(query, langCode).Scan(&langId)

	if langId == 0 {
		langId = 1
	}
	return langId, nil
}

func (r *VerPostgres) GetLevelId(level string) (int, error) {
	var id int

	if level == "" {
		level = "Advanced"
	}

	query := fmt.Sprintf("SELECT id FROM %s WHERE level = $1", levelTable)
	err := r.db.QueryRow(query, level).Scan(&id)
	if err != nil {
		return 0, errors.New("Error query level not found ")
	}

	if id == 0 {
		return 0, errors.New("level not found")
	}

	return id, nil
}

func (r *VerPostgres) GetComments(commentId, userId int) (bool, error) {
	var comment int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = $1 AND user_id = $2")
	err := r.db.QueryRow(query, commentId, userId).Scan(&comment)
	if err != nil {
		return false, err
	}

	if comment > 0 {
		return false, fmt.Errorf("comment %d doesn't exists", comment)
	}

	return true, nil
}
