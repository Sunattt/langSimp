package handler

import (
	"encoding/json"
	"lang/pkg/models"
	"lang/pkg/utils"
	"net/http"
)

// @Summary SignUP
// @Tags auth
// @Description create account
// @ID create account
// @Accept json
// @Produce json
// @Param input body models.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400, 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Failure /auth/sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "invalid payload", err, h.logger)
		return
	}

	id, err := h.service.Authorization.CreateUser(input)
	if err.Code != 0 {
		utils.NewResponseError(w, err.Code, err.Message, err.Error, h.logger)
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	utils.ResponseServer(response, w, h.logger)
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SighIn
// @Tags Auth
// @Description
// @ID login
// @Accept json
// @Produce json
// @Param input body sighInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400, 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.NewResponseError(w, http.StatusBadRequest, "Invalid request payload", err, h.logger)
		return
	}

	token, err := h.service.Authorization.GenerationToken(input.Username, input.Password)
	if err.Code != 0 {
		utils.NewResponseError(w, err.Code, err.Message, err.Error, h.logger)
		return
	}

	response := map[string]interface{}{
		"token": token,
	}

	utils.ResponseServer(response, w, h.logger)
}
