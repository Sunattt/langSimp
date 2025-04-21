package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"lang/pkg/models"
	"lang/pkg/utils"
	"log"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	adminCtx            = "admin"
	langCodeCtx         = "lang_code"
	levelIdCtx          = "level_id"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			utils.NewResponseError(w, http.StatusUnauthorized, "empty auth header", nil, h.logger)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) > 2 {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid auth header", nil, h.logger)
			return
		}

		//parse token
		userId, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid token", err, h.logger)
			return
		}

		// Установка userId в контекст запроса
		ctx := context.WithValue(r.Context(), userCtx, userId)
		r = r.WithContext(ctx)
		// Логирование успешной аутентификации
		log.Printf("User ID %d authenticated successfully.", userId)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) adminIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			utils.NewResponseError(w, http.StatusUnauthorized, "empty auth header", nil, h.logger)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) > 2 {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid auth header", nil, h.logger)
			return
		}

		//parse token
		userId, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid token", err, h.logger)
			return
		}

		admin, err := h.service.IsAdmin(userId)
		if err != nil {
			utils.NewResponseError(w, http.StatusForbidden, "", err, h.logger)
			return
		}

		if !admin {
			err := errors.New("Access denied! ")
			utils.NewResponseError(w, http.StatusForbidden, "", err, h.logger)
			return
		}

		// Установка userId в контекст запроса
		ctx := context.WithValue(r.Context(), adminCtx, admin)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) moderIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)

		if header == "" {
			utils.NewResponseError(w, http.StatusUnauthorized, "empty auth header", nil, h.logger)
			return
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) > 2 {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid auth header", nil, h.logger)
			return
		}

		userId, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid token", err, h.logger)
			return
		}

		moderator, err := h.service.IsModerator(userId)
		if err != nil {
			utils.NewResponseError(w, http.StatusForbidden, "", err, h.logger)
			return
		}

		if !moderator {
			err := errors.New("Access denied! ")
			utils.NewResponseError(w, http.StatusForbidden, "", err, h.logger)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, moderator)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func (h *Handler) langMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		lang := vars[langCodeCtx]

		langId, err := h.service.ValidLangCode(lang)
		if err != nil {
			utils.NewResponseError(w, http.StatusForbidden, "language id not found", err, h.logger)
			return
		}

		ctx := context.WithValue(r.Context(), "lang_id", langId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) getUserId(w http.ResponseWriter, r *http.Request) int {

	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		utils.NewResponseError(w, http.StatusUnauthorized, "user id not found", nil, h.logger)
		return 0
	}

	if id <= 0 {
		utils.NewResponseError(
			w,
			http.StatusBadRequest,
			"Content ID must be positive",
			nil,
			h.logger,
		)
		return 0
	}
	return id
}

func (h *Handler) GetLanId(w http.ResponseWriter, r *http.Request) (int, models.ErrorResponse) {
	// Получения id lang
	langValue := r.Context().Value("lang_id")
	if langValue == nil {
		return 0, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Language identifier missing",
			Error:   errors.New("language ID not found in request context"),
		}
	}

	langID, ok := langValue.(int)
	if !ok {
		return 0, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Invalid language identifier format",
			Error:   errors.New("expected int for lang_id invalid language ID type"),
		}
	}

	if langID <= 0 {
		return 0, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid language identifier",
			Error:   fmt.Errorf("invalid language ID value: %d", langID),
		}
	}

	return langID, models.ErrorResponse{}
}
