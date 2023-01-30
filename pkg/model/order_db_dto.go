package model

type OrderDbDto struct {
	OrderUid          string `db:"order_uid"`
	TrackNumber       string `db:"track_number"`
	Entry             string `db:"entry"`
	DeliveryId        int    `db:"delivery_id"`
	PaymentId         int    `db:"payment_id"`
	Locale            string `db:"locale"`
	InternalSignature string `db:"internal_signature"`
	CustomerId        string `db:"customer_id"`
	DeliveryService   string `db:"delivery_service"`
	Shardkey          string `db:"shardkey"`
	SmId              int    `db:"sm_id"`
	DateCreated       string `db:"date_created"`
	OofShard          string `db:"oof_shard"`
}
