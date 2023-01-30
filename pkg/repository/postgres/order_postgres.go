package postgres

import (
	"L0task/pkg/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) SetOrder(order model.Order, deliveryId int, paymentId int) (string, error) {
	var orderUid string
	dateCreated, conversionErr := time.Parse(time.RFC3339, order.DateCreated)
	if conversionErr != nil {
		fmt.Printf("ERROR PARSING: %s", order.DateCreated)
		return "", conversionErr
	}

	query := fmt.Sprintf("INSERT INTO \"%s\" (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING order_uid", ordersTable)

	row := r.db.QueryRow(query, order.OrderUid, order.TrackNumber, order.Entry, deliveryId, paymentId, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, dateCreated, order.OofShard)
	if err := row.Scan(&orderUid); err != nil {
		return "", err
	}

	return orderUid, nil
}

func (r *OrderPostgres) GetOrderByUid(orderUid string) (model.OrderDbDto, error) {
	var orderDbDto model.OrderDbDto
	query := fmt.Sprintf("SELECT * FROM \"%s\" WHERE order_uid=$1", ordersTable)
	err := r.db.Get(&orderDbDto, query, orderUid)

	return orderDbDto, err
}

func (r *OrderPostgres) GetAllOrders() ([]model.OrderDbDto, error) {
	var orderDbDtos []model.OrderDbDto
	query := fmt.Sprintf("SELECT * FROM \"%s\"", ordersTable)
	err := r.db.Select(&orderDbDtos, query)

	return orderDbDtos, err
}
