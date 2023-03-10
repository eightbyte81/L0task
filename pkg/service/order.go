package service

import (
	"L0task/pkg/model"
	"L0task/pkg/repository"
	"fmt"
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
	if _, dbOrderErr := s.GetOrderByUid(order.OrderUid); dbOrderErr == nil {
		return "", fmt.Errorf("failed to create order: order with uid %s is already exist", order.OrderUid)
	}

	deliveryId, deliveryErr := s.deliveryRepo.SetDelivery(order.Delivery)
	if deliveryErr != nil {
		err := s.RollbackOrderTransaction(deliveryId, 0, model.Order{})
		if err != nil {
			return "", err
		}

		return string(rune(deliveryId)), deliveryErr
	}

	paymentId, paymentErr := s.paymentRepo.SetPayment(order.Payment)
	if paymentErr != nil {
		err := s.RollbackOrderTransaction(deliveryId, paymentId, model.Order{})
		if err != nil {
			return "", err
		}

		return string(rune(paymentId)), paymentErr
	}

	for _, item := range order.Items {
		if itemId, err := s.itemRepo.SetItem(item); err != nil {
			rollErr := s.RollbackOrderTransaction(deliveryId, paymentId, order)
			if rollErr != nil {
				return "", rollErr
			}
			return string(rune(itemId)), err
		}
	}

	orderUid, orderErr := s.repo.SetOrder(order, deliveryId, paymentId)
	if orderErr != nil {
		err := s.RollbackOrderTransaction(deliveryId, paymentId, order)
		if err != nil {
			return "", err
		}

		return orderUid, orderErr
	}

	if orderItemsId, orderItemsErr := s.orderItemsRepo.SetOrderItems(orderUid, order.Items); orderItemsErr != nil {
		err := s.RollbackOrderTransaction(deliveryId, paymentId, order)
		if err != nil {
			return "", err
		}

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

func (s *OrderService) DeleteOrder(orderUid string) error {
	return s.repo.DeleteOrder(orderUid)
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

func (s *OrderService) RollbackOrderTransaction(deliveryId int, paymentId int, order model.Order) error {
	if _, deliveryErr := s.deliveryRepo.GetDeliveryById(deliveryId); deliveryErr == nil {
		deliveryErr = s.deliveryRepo.DeleteDelivery(deliveryId)
		if deliveryErr != nil {
			return deliveryErr
		}
	}

	if _, paymentErr := s.paymentRepo.GetPaymentById(paymentId); paymentErr == nil {
		paymentErr = s.paymentRepo.DeletePayment(paymentId)
		if paymentErr != nil {
			return paymentErr
		}
	}

	for _, item := range order.Items {
		if _, itemErr := s.itemRepo.GetItemById(item.ChrtId); itemErr == nil {
			itemErr = s.itemRepo.DeleteItem(item.ChrtId)
			if itemErr != nil {
				return itemErr
			}
		}
	}

	if _, orderErr := s.GetOrderByUid(order.OrderUid); orderErr == nil {
		orderErr = s.DeleteOrder(order.OrderUid)
		if orderErr != nil {
			return orderErr
		}
	}

	if _, orderItemsErr := s.orderItemsRepo.GetOrderItemsByOrderUid(order.OrderUid); orderItemsErr == nil {
		orderItemsErr = s.orderItemsRepo.DeleteOrderItems(order.OrderUid)
		if orderItemsErr != nil {
			return orderItemsErr
		}
	}

	return nil
}
