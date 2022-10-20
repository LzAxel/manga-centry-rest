package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(ctx *gin.Context, statusCode int, errMessage string) {
	logrus.Error(errMessage)
	ctx.AbortWithStatusJSON(statusCode, errorResponse{Error: errMessage})
}
