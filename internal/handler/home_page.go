package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) HomePage(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"home": "this page is for front-end",
	})
}
