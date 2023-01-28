package service

import (
	"L0task/pkg/model"
	"L0task/pkg/repository"
)

type OrderService struct {
	repo           repository.Order
	deliveryRepo   repository.Delivery
	paymentRepo    repository.Payment
	itemRepo       repository.Item
	orderItemsRepo repository.OrderItems
}

func NewOrderService(repo repository.Order, deliveryRepo repository.Delivery, paymentRepo repository.Payment, itemRepo repository.Item, orderItemsRepo repository.OrderItems) *OrderService {
	return &OrderService{repo: repo, deliveryRepo: deliveryRepo, paymentRepo: paymentRepo, itemRepo: itemRepo, orderItemsRepo: orderItemsRepo}
}

func (s *OrderService) SetOrder(order model.Order) (int, error) {
	deliveryId, deliveryErr := s.deliveryRepo.SetDelivery(order.Delivery)
	if deliveryErr != nil {
		return deliveryId, deliveryErr
	}

	paymentId, paymentErr := s.paymentRepo.SetPayment(order.Payment)
	if paymentErr != nil {
		return paymentId, paymentErr
	}

	for _, item := range order.Items {
		if itemId, err := s.itemRepo.SetItem(item); err != nil {
			return itemId, err
		}
	}

	orderId, orderErr := s.repo.SetOrder(order, deliveryId, paymentId)
	if orderErr != nil {
		return orderId, orderErr
	}

	if orderItemsId, orderItemsErr := s.orderItemsRepo.SetOrderItems(orderId, order.Items); orderItemsErr != nil {
		return orderItemsId, orderItemsErr
	}

	return orderId, nil
}

func (s *OrderService) GetOrderById(orderId int) (model.Order, error) {
	orderDbDto, orderDbDtoErr := s.repo.GetOrderById(orderId)
	if orderDbDtoErr != nil {
		return model.Order{}, orderDbDtoErr
	}

	return s.BuildOrder(orderDbDto)
}

func (s *OrderService) GetAllOrders() ([]model.Order, error) {
	orderDbDtos, orderDbDtosErr := s.repo.GetAllOrders()
	if orderDbDtosErr != nil {
		return []model.Order{}, orderDbDtosErr
	}

	orders := make([]model.Order, len(orderDbDtos))
	for i, orderDbDto := range orderDbDtos {
		if builtOrder, orderBuildingErr := s.BuildOrder(orderDbDto); orderBuildingErr != nil {
			return []model.Order{}, orderBuildingErr
		} else {
			orders[i] = builtOrder
		}
	}

	return orders, nil
}

func (s *OrderService) BuildOrder(orderDbDto model.OrderDbDto) (model.Order, error) {
	delivery, deliveryErr := s.deliveryRepo.GetDeliveryById(orderDbDto.DeliveryId)
	if deliveryErr != nil {
		return model.Order{}, deliveryErr
	}

	payment, paymentErr := s.paymentRepo.GetPaymentById(orderDbDto.PaymentId)
	if paymentErr != nil {
		return model.Order{}, paymentErr
	}

	orderItems, orderItemsErr := s.orderItemsRepo.GetOrderItemsByOrderId(orderDbDto.OrderId)
	if orderItemsErr != nil {
		return model.Order{}, orderItemsErr
	}

	items := make([]model.Item, len(orderItems))
	for i, orderItem := range orderItems {
		if item, itemErr := s.itemRepo.GetItemById(orderItem.ChrtId); itemErr != nil {
			return model.Order{}, itemErr
		} else {
			items[i] = item
		}
	}

	return model.Order{
		OrderUid:          orderDbDto.OrderUid,
		TrackNumber:       orderDbDto.TrackNumber,
		Entry:             orderDbDto.Entry,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            orderDbDto.Locale,
		InternalSignature: orderDbDto.InternalSignature,
		CustomerId:        orderDbDto.CustomerId,
		DeliveryService:   orderDbDto.DeliveryService,
		Shardkey:          orderDbDto.Shardkey,
		SmId:              orderDbDto.SmId,
		DateCreated:       orderDbDto.DateCreated,
		OofShard:          orderDbDto.OofShard,
	}, nil
}
