package models

import "time"

// models chapters _. _._. _ . _._ ._ ._._._. _. _._. _._._. _ ._._ ._._._ . _._.  _._ ._ ._._____
type Chapter struct {
	Id           int
	Title        string
	Description  string
	Priority     int
	Image        string
	CountArticle int
	Created_at   time.Time
	Updated_at   time.Time
	Deleted_at   time.Time
}

type SavedChapter struct {
	Id        int
	ProfileId int
	ChapterId int
}

// models topics __________________________________________________________________________________
type Topics struct {
	Id           int
	Title        string
	Description  string
	Image        string
	Priority     int
	ChapterID    int
	Beginner     bool
	Intermediate bool
	Advanced     bool
	Created_at   time.Time
	Updated_at   time.Time
	Deleted_at   time.Time
}

type SavedTopic struct {
	Id        int
	ProfileId int
	TopicId   int
}

// models article _________________________________________________________________________________
type Article struct {
	Id          int
	Title       string
	Description string
	Image       string
	Comments    string
	Created_at  time.Time
	Updated_at  time.Time
	Deleted_at  time.Time
}
