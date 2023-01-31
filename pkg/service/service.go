package service

import (
	"L0task/pkg/model"
	"L0task/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=../test/mocks/mock.go

type Order interface {
	SetOrder(order model.Order) (string, error)
	SetOrderInCache(order model.Order) error
	SetOrdersFromDbToCache() error
	GetOrderByUid(orderUid string) (model.Order, error)
	GetCachedOrderByUid(orderUid string) (model.Order, error)
	GetAllOrders() ([]model.Order, error)
	GetAllCachedOrders() ([]model.Order, error)
	DeleteOrder(orderUid string) error
	BuildOrder(orderDbDto model.OrderDbDto) (model.Order, error)
	RollbackOrderTransaction(deliveryId int, paymentId int, order model.Order) error
}

type Delivery interface {
	SetDelivery(delivery model.Delivery) (int, error)
	GetDeliveryById(deliveryId int) (model.Delivery, error)
	GetAllDeliveries() ([]model.Delivery, error)
	DeleteDelivery(deliveryId int) error
}

type Payment interface {
	SetPayment(payment model.Payment) (int, error)
	GetPaymentById(paymentId int) (model.Payment, error)
	GetAllPayments() ([]model.Payment, error)
	DeletePayment(paymentId int) error
}

type Item interface {
	SetItem(item model.Item) (int, error)
	GetItemById(itemId int) (model.Item, error)
	GetAllItems() ([]model.Item, error)
	DeleteItem(itemId int) error
}

type OrderItems interface {
	SetOrderItems(orderUid string, items []model.Item) (int, error)
	GetOrderItemsByOrderUid(orderUid string) ([]model.OrderItems, error)
	DeleteOrderItems(orderUid string) error
}

type IService interface {
	Order
	Delivery
	Payment
	Item
	OrderItems
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
