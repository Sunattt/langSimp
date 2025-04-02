package handler

import (
	"encoding/json"
	"lang/pkg/models"
	"lang/pkg/utils"
	"net/http"
)

func (h *Handler) createArticle(w http.ResponseWriter, r *http.Request) {

	var input models.Article

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (h *Handler) getAllArticles(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getArticleById(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) updateArticle(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteArticle(w http.ResponseWriter, r *http.Request) {

}
