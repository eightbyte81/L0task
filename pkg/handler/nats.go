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
		logrus.Fatalf("Failed to connect to nats-streaming server: %s", connErr.Error())
		return stanConn, connErr
	}
	logrus.Print("Connection to nats-streaming server succeed")

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
			return
		}

		cacheErr := n.service.SetOrderInCache(order)
		if cacheErr != nil {
			return
		}

		logrus.Printf("Successfully saved order with uid %s", order.OrderUid)
	})

	if subErr != nil {
		logrus.Fatalf("Failed to subscribe to nats-streaming subject: %s", subErr.Error())
		return subErr
	}

	for {
		if !sub.IsValid() {
			waitGroup.Done()
			break
		}
	}

	if unsubErr := sub.Unsubscribe(); unsubErr != nil {
		logrus.Errorf("Failed to unsubscribe from nats-streaming subject: %s", unsubErr.Error())
		return unsubErr
	}

	logrus.Print("Successfully unsubscribed to nats-streaming subject")

	return nil
}

func (n *Nats) UnmarshalMessage(message *stan.Msg) (model.Order, error) {
	var order model.Order
	jsonErr := json.Unmarshal(message.Data, &order)
	if jsonErr != nil {
		logrus.Errorf("Failed to unmarshal message: %s", jsonErr.Error())
		return order, jsonErr
	}

	validationErr := n.validator.Struct(&order)
	if validationErr != nil {
		logrus.Errorf("Failed to validate struct: %s", validationErr.Error())
		return order, validationErr
	}

	return order, nil
}
