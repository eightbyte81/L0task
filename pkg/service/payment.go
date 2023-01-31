package service

import (
	"L0task/pkg/model"
	"L0task/pkg/repository"
)

type PaymentService struct {
	repo repository.Payment
}

func NewPaymentService(repo repository.Payment) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) SetPayment(payment model.Payment) (int, error) {
	return s.repo.SetPayment(payment)
}

func (s *PaymentService) GetPaymentById(paymentId int) (model.Payment, error) {
	return s.repo.GetPaymentById(paymentId)
}
func (s *PaymentService) GetAllPayments() ([]model.Payment, error) {
	return s.repo.GetAllPayments()
}

func (s *PaymentService) DeletePayment(paymentId int) error {
	return s.repo.DeletePayment(paymentId)
}
