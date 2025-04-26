package service

import (
	"errors"
	"lang/internal/repository"
	"lang/pkg/models"
)

type ProfileSer struct {
	repo repository.ProfileRep
}

func NewProfileSer(repo repository.ProfileRep) *ProfileSer {
	return &ProfileSer{repo: repo}
}

func (s *ProfileSer) SaveChapter(chapterId, userId int) error {

	if chapterId <= 0 || userId <= 0 {
		return errors.New("chapterId or userId can not 0")
	}

	return s.repo.SaveChapter(chapterId, userId)

}

func (s *ProfileSer) SaveArticle(articleId, userId int) error {

	if articleId <= 0 || userId <= 0 {
		return errors.New("articleId or userId can not 0")
	}

	return s.repo.SaveArticle(articleId, userId)
}

func (s *ProfileSer) SaveWord(wordId, userId int) error {
	return s.repo.SaveWord(wordId, userId)
}

func (s *ProfileSer) RemoveSavedChapter(userId, chapterId int) error {
	if chapterId <= 0 || userId <= 0 {
		return errors.New("chapterId or userId can not 0")
	}

	return s.repo.RemoveSavedChapter(userId, chapterId)
}

func (s *ProfileSer) RemoveSavedArticle(userId, articleId int) error {
	return s.repo.RemoveSavedArticle(userId, articleId)
}

func (s *ProfileSer) GetSavedChapters(userId int) ([]models.Chapter, error) {
	return s.repo.GetSavedChapters(userId)
}

func (s *ProfileSer) GetSavedArticles(userId int) ([]models.Article, error) {
	return s.repo.GetSavedArticles(userId)
}
