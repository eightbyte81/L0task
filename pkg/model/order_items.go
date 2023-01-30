package model

type OrderItems struct {
	OrderItemsId int    `db:"order_items_id"`
	OrderUid     string `db:"order_uid"`
	ChrtId       int    `db:"chrt_id"`
}
