package postgres

import (
	"L0task/pkg/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DeliveryPostgres struct {
	db *sqlx.DB
}

func NewDeliveryPostgres(db *sqlx.DB) *DeliveryPostgres {
	return &DeliveryPostgres{db: db}
}

func (r *DeliveryPostgres) SetDelivery(delivery model.Delivery) (int, error) {
	var deliveryId int
	query := fmt.Sprintf("INSERT INTO %s (name, phone, zip, city, address, region, email) values ($1, $2, $3, $4, $5, $6, $7) RETURNING delivery_id", deliveriesTable)

	row := r.db.QueryRow(query, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)
	if err := row.Scan(&deliveryId); err != nil {
		return 0, err
	}

	return deliveryId, nil
}

func (r *DeliveryPostgres) GetDeliveryById(deliveryId int) (model.Delivery, error) {
	var delivery model.Delivery
	query := fmt.Sprintf("SELECT * FROM %s WHERE delivery_id=$1", deliveriesTable)
	err := r.db.Get(&delivery, query, deliveryId)

	return delivery, err
}

func (r *DeliveryPostgres) GetAllDeliveries() ([]model.Delivery, error) {
	var deliveries []model.Delivery
	query := fmt.Sprintf("SELECT * FROM %s", deliveriesTable)
	err := r.db.Get(&deliveries, query)

	return deliveries, err
}
