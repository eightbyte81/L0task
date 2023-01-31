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

func (s *OrderItemsService) SetOrderItems(orderUid string, items []model.Item) (int, error) {
	return s.repo.SetOrderItems(orderUid, items)
}

func (s *OrderItemsService) GetOrderItemsByOrderUid(orderUid string) ([]model.OrderItems, error) {
	return s.repo.GetOrderItemsByOrderUid(orderUid)
}

func (s *OrderItemsService) DeleteOrderItems(orderUid string) error {
	return s.repo.DeleteOrderItems(orderUid)
}
