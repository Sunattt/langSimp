package models

import "time"

type DailyStat struct {
	Date         time.Time `json:"date"`
	MinutesSpent int       `json:"minutes_spent"`
	GoalAchieved bool      `json:"goal_achieved"`
	Goal         int       `json:"daily_goal"`
}

type MonthlyStat struct {
	Year          int     `json:"year"`
	Month         int     `json:"month"`
	TotalMinutes  int     `json:"total_minutes"`
	DaysActivated int     `json:"days_activated"`
	AvgDaily      float64 `json:"average_daily_minutes"`
}

type UserRating struct {
	UserID        int         `json:"user_id"`
	Username      string      `json:"username"`
	TotalMinutes  int         `json:"total_minutes"`
	CurrentStreak int         `json:"current_streak"`
	DailyGoal     int         `json:"daily_goal"`
	Last30Days    []DailyStat `json:"last_30_days"`
}

type ActivitySession struct {
	UserID    int       `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration_minutes"`
}
