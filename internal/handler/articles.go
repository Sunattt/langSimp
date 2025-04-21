package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"lang/pkg/models"
	"lang/pkg/utils"
	"net/http"
	"strconv"
)

func (h *Handler) createArticle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	chapterId, err := strconv.Atoi(vars["chapter_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Invalid chapter ID format", err, h.logger)
		return
	}

	var input *models.Article

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "invalid required data", err, h.logger)
		return
	}

	input.ChapterID = chapterId

	id, errResp := h.service.ArticlePost.CreateArticle(input)
	if errResp.Code != 0 {
		utils.NewResponseError(w, errResp.Code, errResp.Message, errResp.Error, h.logger)
		return
	}

	utils.ResponseServer(map[string]int{"id": id}, w, h.logger)
}

func (h *Handler) getAllArticles(w http.ResponseWriter, r *http.Request) {
	chapterId, err := strconv.Atoi(mux.Vars(r)["chapter_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Invalid chapter ID format", err, h.logger)
		return
	}

	articles, err := h.service.ArticlePost.GetAllArticles(chapterId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to retrieve article", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]interface{}{"data": articles}, w, h.logger)
}

func (h *Handler) getArticleById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["article_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Invalid article ID format", err, h.logger)
		return
	}

	article, err := h.service.ArticlePost.GetArticleById(id)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to retrieve chapter", err, h.logger)
		return
	}

	utils.ResponseServer(article, w, h.logger)
}

func (h *Handler) updateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId, err := strconv.Atoi(vars["article_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Invalid chapter ID format", err, h.logger)
		return
	}

	var input models.UpdateArticle
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "invalid required data", err, h.logger)
		return
	}

	err = h.service.ArticlePost.Update(articleId, input)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to update article", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)
}

func (h *Handler) deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId, err := strconv.Atoi(vars["article_id"])
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Invalid article ID format", err, h.logger)
		return
	}

	err = h.service.ArticlePost.Delete(articleId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to delete article", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)
}
