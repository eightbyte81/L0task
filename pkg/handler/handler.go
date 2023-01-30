package handler

import (
	_ "L0task/docs"
	"L0task/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
