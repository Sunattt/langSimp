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

func (s *VerService) IsAdmin(userId int) (bool, error) {
	return s.repo.IsAdmin(userId)
}
