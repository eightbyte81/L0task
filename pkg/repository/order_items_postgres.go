package repository

import (
	"L0task/pkg/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type OrderItemsPostgres struct {
	db *sqlx.DB
}

func NewOrderItemsPostgres(db *sqlx.DB) *OrderItemsPostgres {
	return &OrderItemsPostgres{db: db}
}

func (r *OrderItemsPostgres) SetOrderItems(orderId int, items []model.Item) (int, error) {
	var lastOrderItemId int

	for _, item := range items {
		query := fmt.Sprintf("INSERT INTO %s (order_id, chrt_id) values ($1, $2) RETURNING chrt_id", orderItemsTable)

		row := r.db.QueryRow(query, orderId, item.ChrtId)
		if err := row.Scan(&item.ChrtId); err != nil {
			return 0, err
		}

		lastOrderItemId = item.ChrtId
	}

	return lastOrderItemId, nil
}

func (r *OrderItemsPostgres) GetOrderItemsByOrderId(orderId int) ([]model.OrderItems, error) {
	var orderItems []model.OrderItems
	query := fmt.Sprintf("SELECT * FROM %s WHERE order_id=$1", orderItemsTable)
	err := r.db.Select(&orderItems, query, orderId)

	return orderItems, err
}
