package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"lang/pkg/helper"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {

	header := c.GetHeader(authorizationHeader)
	if header == "" {
		helper.NewResponseError(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) > 2 {
		helper.NewResponseError(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	//parse token
	userId, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		helper.NewResponseError(c, http.StatusUnauthorized, "invalid token")
		return
	}

	c.Set(userCtx, userId)
}

func (h *Handler) getUserId(c *gin.Context) (int, error) {

	id, ok := c.Get(userCtx)
	if !ok {
		helper.NewResponseError(c, http.StatusUnauthorized, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		helper.NewResponseError(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil

}
