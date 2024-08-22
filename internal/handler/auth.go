package handler

import (
	"lang/pkg/helper"
	"lang/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) sighUp(c *gin.Context) {

	var input models.User

	if err := c.BindJSON(&input); err != nil {
		helper.NewResponseError(c, http.StatusBadRequest, "body is empty!!")
		return
	}

	if err := h.service.Authorization.UserDataValidation(input); err != nil {
		helper.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		helper.NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

type sighInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) sighIn(c *gin.Context) {

	var input sighInInput

	if err := c.BindJSON(&input); err != nil {
		helper.NewResponseError(c, http.StatusBadRequest, "body is empty!")
		return
	}

	token, err := h.service.Authorization.GenerationToken(input.Username, input.Password)
	if err != nil {
		helper.NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
