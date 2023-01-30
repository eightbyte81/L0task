package handler

import (
	"L0task/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	order := router.Group("/")
	{
		order.POST("/", h.setOrder)
		order.GET("/", h.getAllOrders)
		order.GET("/cache/", h.getAllCachedOrders)
		order.GET("/:uid", h.getOrderByUid)
		order.GET("/cache/:uid", h.getCachedOrderByUid)
		order.PUT("/", h.updateOrder)
		order.DELETE("/", h.deleteOrder)
	}

	return router
}
