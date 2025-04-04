package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type VerPostgres struct {
	db *sqlx.DB
}

func NewVerPostgres(db *sqlx.DB) *VerPostgres {
	return &VerPostgres{db: db}
}

func (r *VerPostgres) GetUserActive(userId int, username string) (bool, error) {
	query := fmt.Sprintf("SELECT active FROM %s u WHERE u.username = $1 AND u.id = $2 ", userTable)

	var active bool
	err := r.db.QueryRow(query, username, userId).Scan(&active)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Пользователь не найден
		}
		return false, err // Возвращаем ошибку
	}

	return active, nil // Возвращаем статус активности пользователя
}

func (r *VerPostgres) IsAdmin(userId int) (bool, error) {
	query := fmt.Sprintf("SELECT role FROM %s u WHERE u.id = $1", userTable)

	var roleId int
	err := r.db.QueryRow(query, userId).Scan(&roleId)
	if err != nil {
		return false, err
	}

	if roleId != 1 {
		return false, nil
	}

	return true, nil
}
