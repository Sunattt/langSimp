package handler

import (
	"github.com/gin-gonic/gin"
	"lang/pkg/helper"
	"lang/pkg/models"
	"net/http"
)

func (h *Handler) createChapter(c *gin.Context) {

	id, _ := c.Get(userCtx)

	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	var input models.Chapter

	if err := c.BindJSON(&input); err != nil {
		helper.NewResponseError(c, http.StatusBadRequest, "body is empty")
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

type getAllChaptersResponse struct {
	Data []models.Chapter `json:"data"`
}

func (h *Handler) getAllChapters(c *gin.Context) {

}

func (h *Handler) getChapterById(c *gin.Context) {

}

func (h *Handler) updateChapter(c *gin.Context) {

}

func (h *Handler) deleteChapter(c *gin.Context) {

}
