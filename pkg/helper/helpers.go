package helper

import (

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResp struct {
	Message string `json:"message"`
}

func NewResponseError(c *gin.Context, statusCode int , message string){
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, ErrorResp{Message: message} )
}