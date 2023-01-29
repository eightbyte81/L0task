package handler

import (
	"L0task/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) setOrder(c *gin.Context) {
	var req model.Order

	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, orderErr := h.services.Order.SetOrder(req)
	if orderErr != nil {
		c.String(http.StatusInternalServerError, orderErr.Error())
		return
	}

	if cacheErr := h.services.Order.SetOrderInCache(req); cacheErr != nil {
		c.String(http.StatusInternalServerError, cacheErr.Error())
		return
	}

	c.String(http.StatusOK, "Order created with id %d", id)
}

func (h *Handler) getAllOrders(c *gin.Context) {
	orders, ordersErr := h.services.GetAllOrders()
	if ordersErr != nil {
		c.String(http.StatusInternalServerError, ordersErr.Error())
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) getAllCachedOrders(c *gin.Context) {
	orders, cacheErr := h.services.GetAllCachedOrders()
	if cacheErr != nil {
		c.String(http.StatusInternalServerError, cacheErr.Error())
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) getOrderById(c *gin.Context) {
	orderIdParam := c.Param("id")

	orderId, conversionErr := strconv.Atoi(orderIdParam)
	if conversionErr != nil {
		c.String(http.StatusBadRequest, "Error getting orderId value: %s", orderIdParam)
		return
	}

	order, orderErr := h.services.GetOrderById(orderId)
	if orderErr != nil {
		c.String(http.StatusInternalServerError, orderErr.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) getCachedOrderByUid(c *gin.Context) {
	orderUidParam := c.Param("uid")

	order, cacheErr := h.services.GetCachedOrderByUid(orderUidParam)
	if cacheErr != nil {
		c.String(http.StatusInternalServerError, cacheErr.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) updateOrder(c *gin.Context) {}

func (h *Handler) deleteOrder(c *gin.Context) {}
