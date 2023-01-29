package service

import (
	"L0task/pkg/model"
	"L0task/pkg/repository"
)

type Order interface {
	SetOrder(order model.Order) (int, error)
	SetOrderInCache(order model.Order) error
	SetOrdersFromDbToCache() error
	GetOrderById(orderId int) (model.Order, error)
	GetCachedOrderByUid(orderUid string) (model.Order, error)
	GetAllOrders() ([]model.Order, error)
	GetAllCachedOrders() ([]model.Order, error)
	BuildOrder(orderDbDto model.OrderDbDto) (model.Order, error)
}

type Delivery interface {
	SetDelivery(delivery model.Delivery) (int, error)
	GetDeliveryById(deliveryId int) (model.Delivery, error)
	GetAllDeliveries() ([]model.Delivery, error)
}

type Payment interface {
	SetPayment(payment model.Payment) (int, error)
	GetPaymentById(paymentId int) (model.Payment, error)
	GetAllPayments() ([]model.Payment, error)
}

type Item interface {
	SetItem(item model.Item) (int, error)
	GetItemById(itemId int) (model.Item, error)
	GetAllItems() ([]model.Item, error)
}

type OrderItems interface {
	SetOrderItems(orderId int, items []model.Item) (int, error)
	GetOrderItemsByOrderId(orderId int) ([]model.OrderItems, error)
}

type Service struct {
	Order
	Delivery
	Payment
	Item
	OrderItems
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order:      NewOrderService(repos.Order, repos.Delivery, repos.Payment, repos.Item, repos.OrderItems, repos.OrderCache),
		Delivery:   NewDeliveryService(repos.Delivery),
		Payment:    NewPaymentService(repos.Payment),
		Item:       NewItemService(repos.Item),
		OrderItems: NewOrderItemsService(repos.OrderItems),
	}
}
