package handler

import (
	_ "L0task/docs"
	"L0task/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services service.IService
}

func NewHandler(services service.IService) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	order := router.Group("/")
	{
		order.POST("/", h.SetOrder)
		order.GET("/", h.GetAllOrders)
		order.GET("/cache/", h.GetAllCachedOrders)
		order.GET("/:uid", h.GetOrderByUid)
		order.GET("/cache/:uid", h.GetCachedOrderByUid)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
