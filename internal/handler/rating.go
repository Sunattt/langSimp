package handler

import (
	"github.com/gorilla/mux"
	"lang/pkg/utils"
	"net/http"
	"strconv"
)

func (h *Handler) startSession(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	sessionId, err := h.service.StartSession(userId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to starting session", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]int{"sessionId": sessionId}, w, h.logger)
}

func (h *Handler) endSession(w http.ResponseWriter, r *http.Request) {
	sessionId, err := strconv.Atoi(r.URL.Query().Get("session_id"))
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to parse session_id", err, h.logger)
		return
	}

	err = h.service.EndSession(sessionId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to end session", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]string{"status": "ok"}, w, h.logger)

}

func (h *Handler) getUserRating(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	rating, err := h.service.GetUserRating(userId)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to get user rating", err, h.logger)
		return
	}

	utils.ResponseServer(rating, w, h.logger)
}

func (h *Handler) getMonthlyStats(w http.ResponseWriter, r *http.Request) {
	userId := h.getUserId(w, r)

	year, err := strconv.Atoi(mux.Vars(r)["year"])
	if err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid year", err, h.logger)
		return
	}
	month, err := strconv.Atoi(mux.Vars(r)["month"])

	if err != nil || month < 1 || month > 12 {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid month", err, h.logger)
		return
	}

	stats, err := h.service.GetMonthlyStats(userId, year, month)
	if err != nil {
		utils.NewResponseError(w, http.StatusInternalServerError, "Failed to get month stats", err, h.logger)
		return
	}

	utils.ResponseServer(map[string]interface{}{"stats": stats}, w, h.logger)
}
