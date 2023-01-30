package handler

import (
	"L0task/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type allOrdersResponse struct {
	Orders []model.Order `json:"orders"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
