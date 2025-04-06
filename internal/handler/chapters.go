package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"lang/pkg/models"
	"lang/pkg/utils"
	"net/http"
	"strconv"
)

func (h *Handler) createChapter(w http.ResponseWriter, r *http.Request) {
	var input *models.Chapter

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, err.Error())
		utils.BadRequest(w, err, h.logger)
		return
	}

	id, err := h.service.ChapterPost.Create(input)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, err.Error())
		utils.InternalServerError(w, err, h.logger)
		return
	}

	response := map[string]interface{}{
		"id": id,
	}

	utils.ResponseServer(response, w)
}

type getAllChaptersResponse struct {
	Data []models.Chapter `json:"data"`
}

func (h *Handler) getAllChapters(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	langId, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, err.Error())
		utils.BadRequest(w, err, h.logger)
		return
	}

	chapters, err := h.service.ChapterPost.GetALL(langId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, err.Error())
		utils.InternalServerError(w, err, h.logger)
		return
	}

	response := getAllChaptersResponse{
		Data: chapters,
	}

	utils.ResponseServer(response, w)
}

func (h *Handler) getChapterById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	chap, err := h.service.ChapterPost.GetChapterById(id)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseServer(chap, w)
}

func (h *Handler) updateChapter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chapterId, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, err.Error())
		utils.BadRequest(w, err, h.logger)
		return
	}

	var input models.UpdateChapter
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, err.Error())
		utils.BadRequest(w, err, h.logger)
		return
	}

	err = h.service.ChapterPost.Update(chapterId, input)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, err.Error())
		utils.InternalServerError(w, err, h.logger)
		return
	}

	utils.ResponseServer(utils.StatusResponse{Status: "ok"}, w)
}

func (h *Handler) deleteChapter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chapterId, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, err.Error())
		utils.BadRequest(w, err, h.logger)
		return
	}

	err = h.service.ChapterPost.Delete(chapterId)

	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, err.Error())
		utils.InternalServerError(w, err, h.logger)
		return
	}

	utils.ResponseServer(utils.StatusResponse{Status: "ok"}, w)
}
