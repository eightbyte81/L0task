package cache

import (
	"L0task/pkg/model"
	"fmt"
)

type OrderCache struct {
	orderCache *Cache
}

func NewOrderCache(orderCache *Cache) *OrderCache {
	return &OrderCache{orderCache: orderCache}
}

func (r *OrderCache) SetOrderInCache(orderUid string, order model.Order) error {
	r.orderCache.Mutex.Lock()
	defer r.orderCache.Mutex.Unlock()

	r.orderCache.Data[orderUid] = order

	return nil
}

func (r *OrderCache) GetCachedOrderByUid(orderUid string) (model.Order, error) {
	r.orderCache.Mutex.RLock()
	defer r.orderCache.Mutex.RUnlock()

	if cacheData, found := r.orderCache.Data[orderUid]; found {
		if order, ok := cacheData.(model.Order); ok {
			return order, nil
		}

		return model.Order{}, fmt.Errorf("failed to convert cache value with uid %s to order", orderUid)
	}

	return model.Order{}, fmt.Errorf("failed to get value from cache by uid %s", orderUid)
}

func (r *OrderCache) GetAllCachedOrders() ([]model.Order, error) {
	r.orderCache.Mutex.RLock()
	defer r.orderCache.Mutex.RUnlock()

	if len(r.orderCache.Data) == 0 {
		return []model.Order{}, nil
	}

	orders := make([]model.Order, len(r.orderCache.Data), len(r.orderCache.Data))
	index := 0

	for _, cacheData := range r.orderCache.Data {
		if order, ok := cacheData.(model.Order); ok {
			orders[index] = order
			index++
		} else {
			return nil, fmt.Errorf("failed to convert cache value with uid %s to order", order.OrderUid)
		}
	}

	return orders, nil
}
