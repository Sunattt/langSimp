package models

import (
	"errors"
	"time"
)

//_._._._._._._._._._._._._._._._.Model Chapter _._._._._._._._._._._._._._._._._._.

type Chapter struct {
	Id           int    `json:"-" db:"chapter_id"`
	LanguageId   int    `json:"language" db:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Priority     int    `json:"priority" binding:"required"`
	ImageUrl     string `json:"image" binding:"required"`
	ImageAlt     string `json:"image_alt" binding:"required"`
	CountArticle int    `json:"count_article" binding:"required"`
	Active       bool   `json:"active" binding:"required"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Language struct {
	Id   int    `json:"-" db:"language_id"`
	Code string `json:"code" binding:"required"`
	Lang string `json:"lang" binding:"required"`
}

type UpdateChapter struct {
	Title       *string `json:"title" binding:"request"`
	Description *string `json:"description" binding:"request"`
	ImageUrl    *string `json:"image" binding:"required"`
	ImageAlt    *string `json:"image_alt" binding:"required"`
	UpdateAt    time.Time
}

func (chp UpdateChapter) Validate() error {
	if chp.Title == nil && chp.Description == nil && chp.ImageUrl == nil && chp.ImageAlt == nil {
		return errors.New("Update structure hasn't any change ")
	}
	return nil
}

type LevelArticle struct {
	Id    int `json:"-" db:"id"`
	Level int `json:"level" binding:"required"`
}

// ____________________________________________Model Article______________________________________

type Article struct {
	Id          int
	Title       string
	Description string
	ImageUrl    string
	ImageAlt    string
	Priority    int
	Slug        string
	ChapterID   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

//_______________________________________________Model Content__________________________________

type GrammarContent struct {
	Id          int
	ArticleId   int
	LevelId     int
	Explanation string
	Structure   string
	Example     map[string]string
	Tips        string
	Picture     string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GrammarContentExercises struct {
	Id               int
	GrammarContentId int
	Question         string
	QuestionType     string
	Option           map[string]string
	CorrectAnswer    string
	Explanation      string
	Difficulty       int
	Help             string
}

type GrammarComment struct {
	Id               int
	UserId           int
	GrammarContentId int
	Comment          string
	Rating           int
	LikesCount       int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CommentLikes struct {
	Id        int
	CommentId int
	UserId    int
	CreatedAt time.Time
}
