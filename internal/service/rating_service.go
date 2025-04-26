package service

import (
	"lang/internal/repository"
	"lang/pkg/models"
	"time"
)

type RatingService struct {
	repo repository.RatingPost
}

func NewRatingService(repo repository.RatingPost) *RatingService {
	return &RatingService{repo: repo}
}

func (s *RatingService) StartSession(userId int) (int, error) {
	return s.repo.StartSession(userId)
}

func (s *RatingService) EndSession(sessionId int) error {
	startTime, err := s.GetActivityTime(sessionId)

	if err != nil {
		return err
	}

	// Вычисляем продолжительность
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Minutes())

	// Обновляем сессию
	if err := s.repo.EndSession(sessionId, duration); err != nil {
		return err
	}

	// Обновляем ежедневную активность
	date := time.Now().Truncate(24 * time.Hour)
	if err := s.repo.UpdateDailyActivity(sessionId, date, duration); err != nil {
		return err
	}

	// Обновляем ежемесячную статистику
	year, month, _ := date.Date()
	return s.repo.UpdateMonthlyStats(sessionId, year, int(month), duration)
}

func (s *RatingService) GetMonthlyStats(userID int, year, month int) (*models.MonthlyStat, error) {
	return s.repo.GetMonthlyStats(userID, year, month)
}

func (s *RatingService) GetUserRating(userId int) (*models.UserRating, error) {
	return s.repo.GetUserRating(userId)
}

func (s *RatingService) GetActivityTime(sessionId int) (time.Time, error) {
	return s.repo.GetActivityTime(sessionId)
}
