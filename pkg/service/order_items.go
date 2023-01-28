package service

import (
	"L0task/pkg/model"
	"L0task/pkg/repository"
)

type OrderItemsService struct {
	repo repository.OrderItems
}

func NewOrderItemsService(repo repository.OrderItems) *OrderItemsService {
	return &OrderItemsService{repo: repo}
}

func (s *OrderItemsService) SetOrderItems(orderId int, items []model.Item) (int, error) {
	return s.repo.SetOrderItems(orderId, items)
}

func (s *OrderItemsService) GetOrderItemsByOrderId(orderId int) ([]model.OrderItems, error) {
	return s.repo.GetOrderItemsByOrderId(orderId)
}
