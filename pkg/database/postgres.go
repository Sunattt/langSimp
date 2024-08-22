package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ConfigsDb struct {
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
	SLLmode  string
}

func InitPostgres(cfg ConfigsDb) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName, cfg.SLLmode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
