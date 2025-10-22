package pubsub

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type SubscriptionManager struct {
	Topics     map[string]map[string]*SubscriptionRequest
	Consumers  map[string]Consumer
	Client     *redis.Client
	Connection *websocket.Conn
	DB         *gorm.DB
	ConnLock   *sync.Mutex
	TopicLock  *sync.RWMutex
}

func (s *SubscriptionManager) Subscribe(req *SubscriptionRequest) {
	s.TopicLock.Lock()
	defer s.TopicLock.Unlock()

	_, ok := s.Topics[req.Topic()]

	if !ok {
		tmp := make(map[string]*SubscriptionRequest)
		tmp[req.Name] = req
		s.Topics[req.Topic()] = tmp

		switch req.Topic() {
		case "block":
			s.Consumers[req.Topic()] = NewBlockConsumer(s.Client, s.Connection, s.DB, s.ConnLock, s.TopicLock)
		case "transaction":
			s.Consumers[req.Topic()] = NewTransactionConsumer(s.Client, s.Connection, s.DB, s.ConnLock, s.TopicLock)
		case "event":
			s.Consumers[req.Topic()] = NewEventConsumer(s.Client, s.Connection, s.DB, s.ConnLock, s.TopicLock)
		}

		return
	}

	s.Topics[req.Topic()][req.Name] = req
	s.Consumers[req.Topic()].SendData(&SubscriptionResponse{
		Code: 1,
		Msg:  fmt.Sprintf("Subscribed to %s topic", req.Topic()),
	})
}

func (s *SubscriptionManager) Unsubscribe(req *SubscriptionRequest) {
	s.TopicLock.Lock()
	defer s.TopicLock.Unlock()

	_, ok := s.Topics[req.Topic()]

	if !ok {
		return
	}

	delete(s.Topics[req.Topic()], req.Name)

	if len(s.Topics[req.Topic()]) > 0 {
		s.Consumers[req.Topic()].SendData(&SubscriptionResponse{
			Code: 1,
			Msg:  fmt.Sprintf("Unsubscribed from %s topic", req.Topic()),
		})
		return
	}

	s.Consumers[req.Topic()].Unsubscribe()
	delete(s.Consumers, req.Topic())
	delete(s.Topics, req.Topic())
}
