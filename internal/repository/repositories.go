package repository

import (
	"lang/pkg/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	CheckLangId(landId int) (bool, error)
	GetUser(username, password string) (models.User, error)
	IsEmailFree(email string) (bool, error)
}

type Verification interface {
	GetUserActive(userId int, username string) (bool, error)
	IsAdmin(useId int) (bool, error)
	IsModerator(userId int) (bool, error)
	IsEmailFree(email string) (bool, error)
}

type ChapterPost interface {
	Create(chapter *models.Chapter) (int, error)
	GetALL(landId int) ([]models.Chapter, error)
	GetChapterById(chapterId int) (models.Chapter, error)
	Update(chapterId int, chp models.UpdateChapter) error
	Delete(chapterId int) error
}

type ArticlePost interface {
	Create(article *models.Article) (int, error)
	GetALL(landId int) ([]models.Article, error)
	GetChapterById(chapterId int) (models.Article, error)
	Update(chapterId int, chp models.UpdateChapter) error
	Delete(chapterId int) error
}

type Repository struct {
	Authorization
	ChapterPost
	Verification
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		ChapterPost:   NewChapterPostgres(db),
		Verification:  NewVerPostgres(db),
	}
}
