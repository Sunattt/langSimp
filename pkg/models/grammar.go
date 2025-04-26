package models

import (
	"errors"
	"time"
)

//_._._._._._._._._._._._._._._._.Model Chapter _._._._._._._._._._._._._._._._._._.

type Chapter struct {
	Id           int    `json:"chapter_id" binding:"required" db:"chapter_id"`
	LanguageId   int    `json:"language" binding:"required" db:"language_id"`
	Title        string `json:"title" binding:"required" db:"title"`
	Description  string `json:"description" binding:"required" db:"description"`
	Priority     int    `json:"priority" binding:"required"`
	ImageUrl     string `json:"image_url" binding:"required" db:"image_url"`
	ImageAlt     string `json:"image_alt" binding:"required" db:"image_alt"`
	CountArticle int    `json:"count_articles" binding:"required" db:"count_articles"`
	Active       bool   `json:"active" binding:"required" db:"active"`
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
	Id          int    `json:"article_id" binding:"required" db:"article_id"`
	Title       string `json:"title" binding:"required" db:"title"`
	Description string `json:"description" binding:"required" db:"description"`
	ImageUrl    string `json:"image_url" binding:"required" db:"image_url"`
	ImageAlt    string `json:"image_alt" binding:"required" db:"image_alt"`
	Priority    int
	Active      bool      `json:"active" binding:"required" db:"active"`
	LevelId     int       `json:"level_id" binding:"required" db:"level_id"`
	ChapterID   int       `json:"chapter_id" db:"chapter_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" binding:"required" db:"created_at"`
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type UpdateArticle struct {
	Title       *string   `json:"title" binding:"required"`
	Description *string   `json:"description" binding:"required"`
	LevelId     *int      `json:"level_id" binding:"required" db:"level_id"`
	ImageUrl    *string   `json:"image" binding:"required"`
	ImageAlt    *string   `json:"image_alt" binding:"required"`
	UpdateAt    time.Time `json:"update_at" binding:"required"`
}

func (au UpdateArticle) Validate() error {
	if au.Title == nil && au.Description == nil && au.ImageUrl == nil && au.ImageAlt == nil {
		return errors.New("Update structure hasn't any change ")
	}
	return nil
}

//_______________________________________________Model Content__________________________________

type GrammarContent struct {
	Id          int       `json:"grammar_id" binding:"required" db:"id"`
	ArticleId   int       `json:"article_id" binding:"required" db:"article_id"`
	LevelId     int       `json:"level_id" binding:"required" db:"level_id"`
	Explanation string    `json:"explanation" binding:"required" db:"explanation"`
	Structure   string    `json:"structure" binding:"required" db:"structure"`
	Example     any       `json:"examples" binding:"required" db:"examples"`
	Tips        string    `json:"tips" binding:"required" db:"tips"`
	Picture     string    `json:"picture" binding:"required" db:"picture"`
	Active      bool      `json:"active" binding:"required" db:"active"`
	CreatedAt   time.Time `json:"created_at" binding:"required" db:"created_at"`
	UpdatedAt   time.Time
}

type UpdateContentInput struct {
	Explanation *string            `json:"explanation" binding:"required"`
	Structure   *string            `json:"structure" binding:"required"`
	Example     *map[string]string `json:"example" binding:"required"`
	Tips        *string            `json:"tips" binding:"required"`
	Picture     *string            `json:"picture" binding:"required"`
	Active      *bool              `json:"active" binding:"required"`
	UpdateAt    time.Time          `json:"update_at" binding:"required"`
}

func (cu UpdateContentInput) Validate() error {
	if cu.Tips == nil && cu.Picture == nil && cu.Structure == nil {
		return errors.New("Update structure hasn't any change ")
	}
	if cu.Explanation == nil && cu.Example == nil {
		return errors.New("Update structure hasn't any change ")
	}
	return nil
}

type GrammarLevel struct {
	Level string `json:"level_content" binding:"required"`
}

type GrammarContentExercises struct {
	Id               int    `json:"exercise_id" db:"id"`
	GrammarContentId int    `json:"grammar_content_id" binding:"required" db:"grammar_content_id"`
	Question         string `json:"question" binding:"required" db:"question"`
	QuestionType     string `json:"question_type" binding:"required" db:"question_type"`
	Option           any    `json:"option" binding:"required" db:"option"`
	CorrectAnswer    string `json:"correct_answer" binding:"required" db:"correct_answer"`
	Explanation      string `json:"explanation" binding:"required" db:"explanation"`
	Difficulty       int    `json:"difficulty" binding:"required" db:"difficulty"`
	Help             string `json:"help" binding:"required" db:"help"`
	Active           bool   `json:"active" binding:"required" db:"active"`
}

type UserAnswer struct {
	QuizID    int    `json:"quiz_id" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
}

type QuizResponse struct {
	QuizID      int    `json:"quiz_id" binding:"required"`
	IsCorrect   bool   `json:"is_correct" binding:"required"`
	CorrectAns  string `json:"correct_answer" binding:"required"`
	Explanation string `json:"explanation" binding:"required"`
}

type UpdateGrammarExercise struct {
	Question      *string `json:"question" binding:"required" db:"question"`
	QuestionType  *string `json:"question_type" binding:"required" db:"question_type"`
	Option        *any    `json:"option" binding:"required" db:"option"`
	CorrectAnswer *string `json:"correct_answer" binding:"required" db:"correct_answer"`
	Explanation   *string `json:"explanation" binding:"required" db:"explanation"`
	Difficulty    *int    `json:"difficulty" binding:"required" db:"difficulty"`
	Help          *string `json:"help" binding:"required" db:"help"`
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

type UpdateGrammarComment struct {
	UserId           *int
	GrammarContentId *int
	Comment          *string
	LikeCount        *int
	UpdatedAt        time.Time
}

type CommentLikes struct {
	Id        int
	CommentId int
	UserId    int
	CreatedAt time.Time
}
