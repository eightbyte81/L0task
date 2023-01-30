package repository

import (
	"L0task/pkg/model"
	"L0task/pkg/repository/cache"
	"L0task/pkg/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Order interface {
	SetOrder(order model.Order, deliveryId int, paymentId int) (string, error)
	GetOrderByUid(orderUid string) (model.OrderDbDto, error)
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
	SetOrderItems(orderUid string, items []model.Item) (int, error)
	GetOrderItemsByOrderUid(orderUid string) ([]model.OrderItems, error)
}

type OrderCache interface {
	SetOrderInCache(orderUid string, order model.Order) error
	GetCachedOrderByUid(orderUid string) (model.Order, error)
	GetAllCachedOrders() ([]model.Order, error)
}

type Repository struct {
	Order
	OrderCache
	Delivery
	Payment
	Item
	OrderItems
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order:      postgres.NewOrderPostgres(db),
		OrderCache: cache.NewOrderCache(cache.NewCache()),
		Delivery:   postgres.NewDeliveryPostgres(db),
		Payment:    postgres.NewPaymentPostgres(db),
		Item:       postgres.NewItemPostgres(db),
		OrderItems: postgres.NewOrderItemsPostgres(db),
	}
}
