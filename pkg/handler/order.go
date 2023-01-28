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

func (h *Handler) updateOrder(c *gin.Context) {}

func (h *Handler) deleteOrder(c *gin.Context) {}
