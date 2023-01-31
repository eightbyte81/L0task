package handler

import (
	"L0task/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetOrder
// @Summary SetOrder
// @Description Set order to server
// @ID set-order
// @Accept json
// @Produce json
// @Param req body model.Order true "JSON order request"
// @Success 200 {object} string
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router / [post]
func (h *Handler) SetOrder(c *gin.Context) {
	var req model.Order

	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	uid, orderErr := h.services.SetOrder(req)
	if orderErr != nil {
		newErrorResponse(c, http.StatusInternalServerError, orderErr.Error())
		return
	}

	if cacheErr := h.services.SetOrderInCache(req); cacheErr != nil {
		newErrorResponse(c, http.StatusInternalServerError, cacheErr.Error())
		return
	}

	c.String(http.StatusOK, "Order created with uid %s", uid)
}

// GetAllOrders
// @Summary GetAllOrders
// @Description Get all orders from server
// @ID get-all-orders
// @Accept json
// @Produce json
// @Success 200 {object} allOrdersResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router / [get]
func (h *Handler) GetAllOrders(c *gin.Context) {
	orders, ordersErr := h.services.GetAllOrders()
	if ordersErr != nil {
		newErrorResponse(c, http.StatusInternalServerError, ordersErr.Error())
		return
	}

	c.JSON(http.StatusOK, allOrdersResponse{Orders: orders})
}

// GetAllCachedOrders
// @Summary GetAllCachedOrders
// @Description Get all cached orders from server
// @ID get-all-cached-orders
// @Accept json
// @Produce json
// @Success 200 {object} allOrdersResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /cache/ [get]
func (h *Handler) GetAllCachedOrders(c *gin.Context) {
	orders, cacheErr := h.services.GetAllCachedOrders()
	if cacheErr != nil {
		c.String(http.StatusInternalServerError, cacheErr.Error())
		return
	}

	c.JSON(http.StatusOK, allOrdersResponse{Orders: orders})
}

// GetOrderByUid
// @Summary GetOrderByUid
// @Description Get order by UID from server
// @ID get-order-by-uid
// @Accept json
// @Produce json
// @Param uid path string true "order_uid"
// @Success 200 {object} model.Order
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /{uid} [get]
func (h *Handler) GetOrderByUid(c *gin.Context) {
	orderUidParam := c.Param("uid")

	order, orderErr := h.services.GetOrderByUid(orderUidParam)
	if orderErr != nil {
		c.String(http.StatusInternalServerError, orderErr.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetCachedOrderByUid
// @Summary GetCachedOrderByUid
// @Description Get cached order by UID from server
// @ID get-cached-order-by-uid
// @Accept json
// @Produce json
// @Param uid path string true "order_uid"
// @Success 200 {object} model.Order
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /cache/{uid} [get]
func (h *Handler) GetCachedOrderByUid(c *gin.Context) {
	orderUidParam := c.Param("uid")

	order, cacheErr := h.services.GetCachedOrderByUid(orderUidParam)
	if cacheErr != nil {
		c.String(http.StatusInternalServerError, cacheErr.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}
