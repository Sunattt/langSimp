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
	ValidLangCode(langCode string) (int, error)
	GetLevelId(level string) (int, error)
	GetComments(commentId, userId int) (bool, error)
}

type ChapterPost interface {
	Create(chapter *models.Chapter) (int, error)
	GetALL(landId int) ([]models.Chapter, error)
	GetChapterById(chapterId, langId int) (models.Chapter, error)
	Update(chapterId int, chp models.UpdateChapter) error
	Delete(chapterId int) error
}

type ArticlePost interface {
	CreateArticle(article *models.Article) (int, error)
	GetAllArticles(chapterId int) ([]models.Article, error)
	GetArticleById(articleId int) (models.Article, error)
	UpdateArticle(articleId int, chp models.UpdateArticle) error
	Delete(articleId int) error
}

type ContentPost interface {
	CreateContent(article *models.GrammarContent) (int, error)
	GetCourseById(articleId, levelId int) (models.GrammarContent, error)
	UpdateContent(contentId, levelId int, input models.UpdateContentInput) error
	DeleteContent(contentId, levelId int) error
	CreateExercise(article *models.GrammarContentExercises) (int, error)
	GetExerciseById(articleId int) (models.GrammarContentExercises, error)
	UpdateExercise(contentId int, input models.UpdateGrammarExercise) error
	DeleteExercise(contentId int) error
	CreateComment(userId, contentId int, input models.GrammarComment) (int, error)
	GetAllComments(contentsId int) ([]models.GrammarComment, error)
	UpdateComment(commentId int, input models.UpdateGrammarComment) error
	DeleteComment(commentId int) error
	AddLike(userId, contentId int) error
	RemoveLike(userId, contentId int) error
	CheckLikeExists(userId, commentId int) (bool, error)
}

type Repository struct {
	Authorization
	ChapterPost
	Verification
	ArticlePost
	ContentPost
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		ChapterPost:   NewChapterPostgres(db),
		Verification:  NewVerPostgres(db),
		ArticlePost:   NewArticlePostgres(db),
		ContentPost:   NewContentPostgres(db),
	}
}
