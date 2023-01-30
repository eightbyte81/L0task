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
	orderCacheRepo repository.OrderCache
}

func NewOrderService(repo repository.Order, deliveryRepo repository.Delivery, paymentRepo repository.Payment, itemRepo repository.Item, orderItemsRepo repository.OrderItems, orderCacheRepo repository.OrderCache) *OrderService {
	return &OrderService{repo: repo, deliveryRepo: deliveryRepo, paymentRepo: paymentRepo, itemRepo: itemRepo, orderItemsRepo: orderItemsRepo, orderCacheRepo: orderCacheRepo}
}

func (s *OrderService) SetOrder(order model.Order) (string, error) {
	deliveryId, deliveryErr := s.deliveryRepo.SetDelivery(order.Delivery)
	if deliveryErr != nil {
		return string(rune(deliveryId)), deliveryErr
	}

	paymentId, paymentErr := s.paymentRepo.SetPayment(order.Payment)
	if paymentErr != nil {
		return string(rune(paymentId)), paymentErr
	}

	for _, item := range order.Items {
		if itemId, err := s.itemRepo.SetItem(item); err != nil {
			return string(rune(itemId)), err
		}
	}

	orderUid, orderErr := s.repo.SetOrder(order, deliveryId, paymentId)
	if orderErr != nil {
		return orderUid, orderErr
	}

	if orderItemsId, orderItemsErr := s.orderItemsRepo.SetOrderItems(orderUid, order.Items); orderItemsErr != nil {
		return string(rune(orderItemsId)), orderItemsErr
	}

	return orderUid, nil
}

func (s *OrderService) SetOrderInCache(order model.Order) error {
	return s.orderCacheRepo.SetOrderInCache(order.OrderUid, order)
}

func (s *OrderService) SetOrdersFromDbToCache() error {
	orders, ordersErr := s.GetAllOrders()
	if ordersErr != nil {
		return ordersErr
	}

	for i := 0; i < len(orders); i++ {
		if cacheErr := s.SetOrderInCache(orders[i]); cacheErr != nil {
			return cacheErr
		}
	}

	return nil
}

func (s *OrderService) GetOrderByUid(orderUid string) (model.Order, error) {
	orderDbDto, orderDbDtoErr := s.repo.GetOrderByUid(orderUid)
	if orderDbDtoErr != nil {
		return model.Order{}, orderDbDtoErr
	}

	return s.BuildOrder(orderDbDto)
}

func (s *OrderService) GetCachedOrderByUid(orderUid string) (model.Order, error) {
	return s.orderCacheRepo.GetCachedOrderByUid(orderUid)
}

func (s *OrderService) GetAllOrders() ([]model.Order, error) {
	orderDbDtos, orderDbDtosErr := s.repo.GetAllOrders()
	if orderDbDtosErr != nil {
		return []model.Order{}, orderDbDtosErr
	}

	orders := make([]model.Order, len(orderDbDtos), len(orderDbDtos))
	for i, orderDbDto := range orderDbDtos {
		if builtOrder, orderBuildingErr := s.BuildOrder(orderDbDto); orderBuildingErr != nil {
			return []model.Order{}, orderBuildingErr
		} else {
			orders[i] = builtOrder
		}
	}

	return orders, nil
}

func (s *OrderService) GetAllCachedOrders() ([]model.Order, error) {
	return s.orderCacheRepo.GetAllCachedOrders()
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

	orderItems, orderItemsErr := s.orderItemsRepo.GetOrderItemsByOrderUid(orderDbDto.OrderUid)
	if orderItemsErr != nil {
		return model.Order{}, orderItemsErr
	}

	items := make([]model.Item, len(orderItems), len(orderItems))
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
