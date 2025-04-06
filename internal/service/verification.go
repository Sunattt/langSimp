package service

import "lang/internal/repository"

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
