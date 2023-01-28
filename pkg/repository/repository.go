package repository

import (
	"L0task/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Order interface {
	SetOrder(order model.Order, deliveryId int, paymentId int) (int, error)
	GetOrderById(orderId int) (model.OrderDbDto, error)
	GetAllOrders() ([]model.OrderDbDto, error)
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

type Repository struct {
	Order
	Delivery
	Payment
	Item
	OrderItems
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order:      NewOrderPostgres(db),
		Delivery:   NewDeliveryPostgres(db),
		Payment:    NewPaymentPostgres(db),
		Item:       NewItemPostgres(db),
		OrderItems: NewOrderItemsPostgres(db),
	}
}
