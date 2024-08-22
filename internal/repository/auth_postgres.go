package repository

import (
	"fmt"
	"lang/pkg/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

const (
	userTable = "users"
)

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int

	user.Created_at = time.Now()

	query := fmt.Sprintf("INSERT INTO %s (username, gender, language, email, password_hash, created_at) VALUE ($1, $2, $3, $4, $5, $6) RETURNING id", userTable)
	row := r.db.QueryRow(query, user.Username, user.Gender, user.Language, user.Email, user.Password, user.Created_at)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1 AND  password_hash = $2", userTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
