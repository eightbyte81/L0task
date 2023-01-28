package service

import (
	"L0task/pkg/model"
	"L0task/pkg/repository"
)

type DeliveryService struct {
	repo repository.Delivery
}

func NewDeliveryService(repo repository.Delivery) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (s *DeliveryService) SetDelivery(delivery model.Delivery) (int, error) {
	return s.repo.SetDelivery(delivery)
}

func (s *DeliveryService) GetDeliveryById(deliveryId int) (model.Delivery, error) {
	return s.repo.GetDeliveryById(deliveryId)
}

func (s *DeliveryService) GetAllDeliveries() ([]model.Delivery, error) {
	return s.repo.GetAllDeliveries()
}
