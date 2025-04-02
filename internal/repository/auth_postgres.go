package repository

import (
	"errors"
	"fmt"
	"lang/pkg/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

const (
	userTable = "users"
)

func (r *AuthPostgres) CheckLangId(lang int) (bool, error) {
	var LangIdExist bool

	if lang > 1 {
		return false, errors.New("Invalid language, doesn't choose ")
	}

	query := fmt.Sprint("SELECT EXISTS(SELECT * FROM languages WHERE language_id = $1);")
	err := r.db.QueryRow(query, lang).Scan(&LangIdExist)
	if err != nil {
		return false, err
	}

	if !LangIdExist {
		return false, errors.New("Language not exist ")
	}

	return LangIdExist, nil
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	user.CreatedAt = time.Now()

	query := fmt.Sprintf("INSERT INTO %s (username, gender, language_id, birthday, email, password_hash, created_at) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id", userTable)
	row := r.db.QueryRow(query, user.Username, user.Gender, user.Language, user.Birthday, user.Email, user.Password, user.CreatedAt)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) IsEmailFree(email string) (bool, error) {
	var answer bool

	query := fmt.Sprint("SELECT EXISTS(SELECT * FROM users WHERE email = $1);")
	err := r.db.QueryRow(query, email).Scan(&answer)
	if err != nil {
		return false, err
	}

	return answer, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1 AND  password_hash = $2", userTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
