package postgres

import (
	"L0task/pkg/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PaymentPostgres struct {
	db *sqlx.DB
}

func NewPaymentPostgres(db *sqlx.DB) *PaymentPostgres {
	return &PaymentPostgres{db: db}
}

func (r *PaymentPostgres) SetPayment(payment model.Payment) (int, error) {
	var paymentId int
	query := fmt.Sprintf("INSERT INTO %s (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING payment_id", paymentsTable)

	row := r.db.QueryRow(query, payment.Transaction, payment.RequestId, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	if err := row.Scan(&paymentId); err != nil {
		return 0, err
	}

	return paymentId, nil
}

func (r *PaymentPostgres) GetPaymentById(paymentId int) (model.Payment, error) {
	var payment model.Payment
	query := fmt.Sprintf("SELECT * FROM %s WHERE payment_id=$1", paymentsTable)
	err := r.db.Get(&payment, query, paymentId)

	return payment, err
}
func (r *PaymentPostgres) GetAllPayments() ([]model.Payment, error) {
	var payments []model.Payment
	query := fmt.Sprintf("SELECT * FROM %s", paymentsTable)
	err := r.db.Select(&payments, query)

	return payments, err
}
