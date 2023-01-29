package postgres

import (
	"L0task/pkg/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ItemPostgres struct {
	db *sqlx.DB
}

func NewItemPostgres(db *sqlx.DB) *ItemPostgres {
	return &ItemPostgres{db: db}
}

func (r *ItemPostgres) SetItem(item model.Item) (int, error) {
	var itemId int
	query := fmt.Sprintf("INSERT INTO %s (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING chrt_id", itemsTable)

	row := r.db.QueryRow(query, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
	if err := row.Scan(&itemId); err != nil {
		return 0, err
	}

	return itemId, nil
}

func (r *ItemPostgres) GetItemById(itemId int) (model.Item, error) {
	var item model.Item
	query := fmt.Sprintf("SELECT * FROM %s WHERE chrt_id=$1", itemsTable)
	err := r.db.Get(&item, query, itemId)

	return item, err
}
func (r *ItemPostgres) GetAllItems() ([]model.Item, error) {
	var items []model.Item
	query := fmt.Sprintf("SELECT * FROM %s", itemsTable)
	err := r.db.Select(&items, query)

	return items, err
}
