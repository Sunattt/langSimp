package models

import "time"

type User struct {
	Id         int    `json:"-" bd:"id"`
	Username   string `json:"username" binding:"required"`
	Gender     string `json:"gender" binding:"required"`
	Language   string `json:"language" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Status     bool
	RoleId     int
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}

type ProfileUser struct {
	Id             int    `json:"-" bd:"id"`
	UserId         int    `json:"userId" binding:"required"`
	GrammarMark    int    `json:"grammarMark" binding:"required"`
	VocabularyMark int    `json:"vocabularyMark" binding:"required"`
	PhotoProfile   string `json:"photoProfile" binding:"required"`
}
