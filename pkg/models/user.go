package models

import "time"

type User struct {
	Id            int    `json:"-" bd:"id"`
	Username      string `json:"username" binding:"required"`
	DailyGoal     int    `json:"daily_goal" db:"daily_goal_minutes"`
	CurrentStreak int    `json:"streak_days" db:"streak_days"`
	Gender        int    `json:"gender" binding:"required"`
	Birthday      string `json:"birthday" binding:"required"`
	Language      int    `json:"language" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
	PhotoUrl      string `json:"photo_url" binding:"required"`
	Active        bool   `json:"active" binding:"required"`
	RoleId        int    `json:"roleId" binding:"required"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

type UserSavedChapter struct {
	Id        int
	ProfileId int
	ChapterId int
}

type UserSavedArticle struct {
	Id        int
	ProfileId int
	ArticleId int
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error
}
