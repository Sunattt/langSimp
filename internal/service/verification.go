package service

import (
	"lang/internal/repository"
)

type VerService struct {
	repo repository.Verification
}

func NewVerService(repo repository.Verification) *VerService {
	return &VerService{repo: repo}
}

func (s *VerService) GetUserActive(userId int, username string) (bool, error) {
	return s.repo.GetUserActive(userId, username)
}

func (s *VerService) IsEmailFree(email string) (bool, error) {
	return s.repo.IsEmailFree(email)
}

func (s *VerService) IsAdmin(userId int) (bool, error) {
	return s.repo.IsAdmin(userId)
}

func (s *VerService) IsModerator(userId int) (bool, error) {
	return s.repo.IsModerator(userId)
}

func (s *VerService) ValidLangCode(langCode string) (int, error) {
	return s.repo.ValidLangCode(langCode)
}

func (s *VerService) GetLevelId(level string) (int, error) {
	return s.repo.GetLevelId(level)
}

func (s *VerService) GetComments(commentId, userId int) (bool, error) {
	return s.repo.GetComments(commentId, userId)
}

func (s *VerService) CheckLikeExists(userId, commentId int) (bool, error) {
	return false, nil
}

func (s *VerService) IsValidChapterId(chapterId int) (bool, error) {
	return s.repo.IsValidChapterId(chapterId)
}

func (s *VerService) IsValidArticleId(articleId int) (bool, error) {
	return s.repo.IsValidArticleId(articleId)
}
