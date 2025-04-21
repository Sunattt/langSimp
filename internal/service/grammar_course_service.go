package service

import (
	"errors"
	"lang/internal/repository"
	"lang/pkg/models"
	"net/http"
	"strings"
)

type CourseService struct {
	repo repository.ContentPost
}

var (
	errContentNotFound  = errors.New("content not found")
	errUserNotFound     = errors.New("user not found")
	errDuplicateComment = errors.New("duplicate comment")
)

func NewCourseService(repo repository.ContentPost) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) CreateContent(input *models.GrammarContent) (int, error) {
	return s.repo.CreateContent(input)
}

func (s *CourseService) GetCourseById(contentId, levelId int) (models.GrammarContent, error) {
	return s.repo.GetCourseById(contentId, levelId)
}

func (s *CourseService) UpdateContent(contentId, levelId int, input models.UpdateContentInput) error {
	return s.repo.UpdateContent(contentId, levelId, input)
}

func (s *CourseService) DeleteContent(contentId, levelId int) error {
	return s.repo.DeleteContent(contentId, levelId)
}

func (s *CourseService) CreateExercise(article *models.GrammarContentExercises) (int, error) {

	return s.repo.CreateExercise(article)
}

func (s *CourseService) GetExerciseById(articleId, levelId int) (models.GrammarContentExercises, error) {
	return s.repo.GetExerciseById(articleId)
}

func (s *CourseService) UpdateExercise(contentId int, input models.UpdateGrammarExercise) error {
	return s.repo.UpdateExercise(contentId, input)
}

func (s *CourseService) DeleteExercise(contentId, levelId int) error {
	return s.repo.DeleteExercise(contentId)
}

func (s *CourseService) CreateComment(userId, content int, input models.GrammarComment) (int, models.ErrorResponse) {

	if strings.TrimSpace(input.Comment) == "" {
		return 0, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Comment text cannot be empty",
		}
	}

	commentId, err := s.repo.CreateComment(userId, content, input)
	if err != nil {
		if errors.Is(err, errContentNotFound) {
			return 0, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Content not found",
				Error:   err,
			}
		}

		if errors.Is(err, errUserNotFound) {
			return 0, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "User not found",
				Error:   err,
			}
		}

		return 0, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create comment",
			Error:   err,
		}
	}

	return commentId, models.ErrorResponse{}

}

func (s *CourseService) GetAllComments(contentId int) ([]models.GrammarComment, error) {
	return s.repo.GetAllComments(contentId)
}

func (s *CourseService) UpdateComment(commentId int, input models.UpdateGrammarComment) error {
	return s.repo.UpdateComment(commentId, input)
}

func (s *CourseService) DeleteComment(commentId int) error {
	return s.repo.DeleteComment(commentId)
}
func (s *CourseService) LikeComment(userId, commentId int) error {
	exists, err := s.repo.CheckLikeExists(userId, commentId)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already liked this comment")
	}

	// Добавляем лайк
	return s.repo.AddLike(userId, commentId)
}

func (s *CourseService) RemoveLike(userId, commentId int) error {
	return s.repo.RemoveLike(userId, commentId)
}
