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

	langId, errResp := h.GetLanId(w, r)
	if errResp.Code != 0 {
		utils.NewResponseError(w, errResp.Code, errResp.Message, errResp.Error, h.logger)
		return
	}

	var input *models.Chapter

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "invalid required data", err, h.logger)
		return
	}

	input.LanguageId = langId

	id, err := h.service.ChapterPost.Create(input)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "failed to create chapter", err, h.logger)
		return
	}

	response := map[string]interface{}{
		"id": id,
	}

	utils.ResponseServer(response, w, h.logger)

}

type getAllChaptersResponse struct {
	Data []models.Chapter `json:"data"`
}

func (h *Handler) getAllChapters(w http.ResponseWriter, r *http.Request) {

	langId, errResp := h.GetLanId(w, r)

	if errResp.Code != 0 {
		utils.NewResponseError(w, errResp.Code, errResp.Message, errResp.Error, h.logger)
		return
	}

	chapters, err := h.service.ChapterPost.GetALL(langId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to retrieve chapters", err, h.logger)
		return
	}

	response := getAllChaptersResponse{
		Data: chapters,
	}

	utils.ResponseServer(response, w, h.logger)

}

func (h *Handler) getChapterById(w http.ResponseWriter, r *http.Request) {
	langId, errResp := h.GetLanId(w, r)
	if errResp.Code != 0 {
		utils.NewResponseError(w, errResp.Code, errResp.Message, errResp.Error, h.logger)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["chapter_id"])

	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid chapter ID format", err, h.logger)
		return
	}

	chap, err := h.service.ChapterPost.GetChapterById(id, langId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to retrieve chapter", err, h.logger)
		return
	}

	utils.ResponseServer(chap, w, h.logger)

}

func (h *Handler) updateChapter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chapterId, err := strconv.Atoi(vars["chapter_id"])

	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid chapter ID format", err, h.logger)
		return
	}

	var input models.UpdateChapter
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Failed to retrieve data", err, h.logger)
		return
	}

	err = h.service.ChapterPost.Update(chapterId, input)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to update chapter", err, h.logger)
		return
	}

	utils.ResponseServer(utils.StatusResponse{Status: "ok"}, w, h.logger)
}

func (h *Handler) deleteChapter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chapterId, err := strconv.Atoi(vars["chapter_id"])

	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid chapter ID format", err, h.logger)
		return
	}

	err = h.service.ChapterPost.Delete(chapterId)

	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to delete chapter", err, h.logger)
		return
	}

	utils.ResponseServer(utils.StatusResponse{Status: "ok"}, w, h.logger)
}
