package handler

import (
	"context"
	"errors"
	"lang/pkg/utils"
	"log"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	adminCtx            = "admin"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			utils.NewResponseError(w, http.StatusUnauthorized, "empty auth header")
			utils.UnauthorizedError(w, errors.New("Empty error "), h.logger)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) > 2 {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid auth header")
			utils.BadRequest(w, errors.New("invalid auth header"), h.logger)
			return
		}

		//parse token
		userId, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid token")
			utils.BadRequest(w, errors.New("invalid token"), h.logger)
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
			utils.NewResponseError(w, http.StatusUnauthorized, "empty auth header")
			utils.UnauthorizedError(w, errors.New("Empty error "), h.logger)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) > 2 {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid auth header")
			utils.BadRequest(w, errors.New("invalid auth header"), h.logger)
			return
		}

		//parse token
		userId, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			utils.NewResponseError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		admin, err := h.service.IsAdmin(userId)
		if err != nil {
			utils.NewResponseError(w, http.StatusForbidden, err.Error())
			utils.Forbidden(w, err, h.logger)
			return
		}

		if !admin {
			err := errors.New("Access denied! ")
			utils.NewResponseError(w, http.StatusForbidden, err.Error())
			utils.Forbidden(w, err, h.logger)
			return
		}

		// Установка userId в контекст запроса
		ctx := context.WithValue(r.Context(), adminCtx, admin)
		r = r.WithContext(ctx)

	})
}

func (h *Handler) getUserId(w http.ResponseWriter, r *http.Request) (int, error) {

	id, ok := r.Context().Value(userCtx).(int)
	if !ok {
		utils.NewResponseError(w, http.StatusUnauthorized, "user id not found")
		utils.UnauthorizedError(w, errors.New("user id not found"), h.logger)
		return 0, errors.New("user id not found")
	}

	return id, nil

}
