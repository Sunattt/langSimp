package service

import (
	"lang/internal/repository"
	"lang/pkg/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, models.ErrorResponse)
	GenerationToken(username, password string) (string, models.ErrorResponse)
	ParseToken(token string) (int, error)
}

type Verification interface {
	GetUserActive(userId int, username string) (bool, error)
	IsAdmin(userId int) (bool, error)
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
	Update(chapterId int, input models.UpdateChapter) error
	Delete(chapterId int) error
}

type ArticlePost interface {
	CreateArticle(article *models.Article) (int, models.ErrorResponse)
	GetAllArticles(chapterId int) ([]models.Article, error)
	GetArticleById(articleId int) (models.Article, error)
	Update(articleId int, chp models.UpdateArticle) error
	Delete(articleId int) error
}

type CoursePost interface {
	CreateContent(article *models.GrammarContent) (int, error)
	GetCourseById(articleId, levelId int) (models.GrammarContent, error)
	UpdateContent(contentId, levelId int, input models.UpdateContentInput) error
	DeleteContent(contentId, levelId int) error
	CreateExercise(article *models.GrammarContentExercises) (int, error)
	GetExerciseById(articleId, levelId int) (models.GrammarContentExercises, error)
	UpdateExercise(contentId int, input models.UpdateGrammarExercise) error
	DeleteExercise(contentId, levelId int) error
	CreateComment(userId, contentId int, comment models.GrammarComment) (int, models.ErrorResponse)
	GetAllComments(contentId int) ([]models.GrammarComment, error)
	UpdateComment(commentId int, input models.UpdateGrammarComment) error
	DeleteComment(commentId int) error
	LikeComment(userId, contentId int) error
	RemoveLike(userId, commentId int) error
}

type Ser struct {
	Authorization
	ChapterPost
	Verification
	CoursePost
	ArticlePost
}

func NewService(repo *repository.Repository) *Ser {
	return &Ser{
		Authorization: NewAuthService(repo.Authorization, repo.Verification),
		ChapterPost:   NewChapterService(repo.ChapterPost),
		Verification:  NewVerService(repo.Verification),
		CoursePost:    NewCourseService(repo.ContentPost),
		ArticlePost:   NewArticleService(repo.ArticlePost),
	}
}
