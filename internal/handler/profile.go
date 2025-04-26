package handler

import (
	"github.com/gorilla/mux"
	"lang/pkg/utils"
	"net/http"
	"strconv"
)

func (h *Handler) saveChapter(w http.ResponseWriter, r *http.Request) {
	contentId, err := strconv.Atoi(mux.Vars(r)["content_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid content id format", err, h.logger)
		return
	}

	userId := h.getUserId(w, r)

	err = h.service.SaveChapter(contentId, userId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to save chapter", err, h.logger)
	}

	utils.ResponseServer(map[string]string{"status": "chapter_saved"}, w, h.logger)
}

func (h *Handler) saveArticle(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	articleId, err := strconv.Atoi(mux.Vars(r)["article_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid article id format", err, h.logger)
		return
	}

	err = h.service.SaveArticle(articleId, userId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to save article", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "article_saved"}, w, h.logger)
}

func (h *Handler) saveWord(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getSavedChapter(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	savedChapters, err := h.service.GetSavedChapters(userId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to save chapters", err, h.logger)
		return
	}

	response := map[string]interface{}{
		"saved_chapters": savedChapters,
	}

	utils.ResponseServer(response, w, h.logger)
}

func (h *Handler) getSavedArticle(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	savedArticles, err := h.service.GetSavedArticles(userId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to save articles", err, h.logger)
		return
	}

	response := map[string]interface{}{
		"saved_articles": savedArticles,
	}

	utils.ResponseServer(response, w, h.logger)
}

func (h *Handler) getSavedWord(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) removeSavedChapter(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	chapterId, err := strconv.Atoi(mux.Vars(r)["chapter_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid chapter id format", err, h.logger)
		return
	}

	err = h.service.RemoveSavedChapter(userId, chapterId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to save chapter", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "chapter_removed"}, w, h.logger)
}

func (h *Handler) removeSavedArticle(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	articleId, err := strconv.Atoi(mux.Vars(r)["article_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid article id format", err, h.logger)
		return
	}

	err = h.service.RemoveSavedArticle(userId, articleId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to save article", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "article_removed"}, w, h.logger)
}

func (h *Handler) removeSavedWord(w http.ResponseWriter, r *http.Request) {}
