package handler

import (
	"L0task/pkg/model"
	"L0task/pkg/service"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"sync"
)

type Nats struct {
	service   *service.Service
	validator *validator.Validate
}

func NewNats(service *service.Service, validator *validator.Validate) *Nats {
	return &Nats{service: service, validator: validator}
}

func (n *Nats) Connect(clusterId, clientId, natsUrl string) (stan.Conn, error) {
	stanConn, connErr := stan.Connect(clusterId, clientId, stan.NatsURL(natsUrl))
	if connErr != nil {
		logrus.Fatalf("[NATS] Failed to connect to nats-streaming server: %s", connErr.Error())
		return stanConn, connErr
	}
	logrus.Print("[NATS] Connection to nats-streaming server succeed")

	return stanConn, nil
}

func (n *Nats) Subscribe(waitGroup *sync.WaitGroup, stanConn stan.Conn, natsSubject string) error {
	defer waitGroup.Done()

	sub, subErr := stanConn.Subscribe(natsSubject, func(message *stan.Msg) {
		order, msgErr := n.UnmarshalMessage(message)
		if msgErr != nil {
			return
		}

		_, dbErr := n.service.SetOrder(order)
		if dbErr != nil {
			logrus.Printf("[NATS] Failed to save order: %s", dbErr.Error())
			return
		}

		cacheErr := n.service.SetOrderInCache(order)
		if cacheErr != nil {
			logrus.Printf("[NATS] Failed to set order wih uid %s in cache: %s", order.OrderUid, cacheErr.Error())
			return
		}

		logrus.Printf("[NATS] Successfully saved order with uid %s", order.OrderUid)
	})

	if subErr != nil {
		logrus.Fatalf("[NATS] Failed to subscribe to nats-streaming subject: %s", subErr.Error())
		return subErr
	}

	for {
		if !sub.IsValid() {
			waitGroup.Done()
			break
		}
	}

	if unsubErr := sub.Unsubscribe(); unsubErr != nil {
		logrus.Errorf("[NATS] Failed to unsubscribe from nats-streaming subject: %s", unsubErr.Error())
		return unsubErr
	}

	logrus.Print("[NATS] Successfully unsubscribed to nats-streaming subject")

	return nil
}

func (n *Nats) UnmarshalMessage(message *stan.Msg) (model.Order, error) {
	var order model.Order
	jsonErr := json.Unmarshal(message.Data, &order)
	if jsonErr != nil {
		logrus.Errorf("[NATS] Failed to unmarshal message: %s", jsonErr.Error())
		return order, jsonErr
	}

	validationErr := n.validator.Struct(&order)
	if validationErr != nil {
		logrus.Errorf("[NATS] Failed to validate struct: %s", validationErr.Error())
		return order, validationErr
	}

	return order, nil
}
