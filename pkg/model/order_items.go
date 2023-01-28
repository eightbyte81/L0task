package model

type OrderItems struct {
	OrderItemsId int `db:"order_items_id"`
	OrderId      int `db:"order_id"`
	ChrtId       int `db:"chrt_id"`
}
