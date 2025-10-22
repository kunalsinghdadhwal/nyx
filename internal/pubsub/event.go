package pubsub

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	d "github.com/kunalsinghdadhwal/nyx/internal/data"
	"github.com/kunalsinghdadhwal/nyx/pkg/logger"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type EventConsumer struct {
	Client     *redis.Client
	Requests   map[string]*SubscriptionRequest
	Connection *websocket.Conn
	Pubsub     *redis.PubSub
	DB         *gorm.DB
	ConnLock   *sync.Mutex
	TopicLock  *sync.RWMutex
}

func (e *EventConsumer) Subscribe() {
	e.Pubsub = e.Client.Subscribe(context.Background(), "event")
}

func (e *EventConsumer) Listen() {
	for {
		msg, err := e.Pubsub.ReceiveTimeout(context.Background(), time.Duration(1)*time.Second)
		if err != nil {
			continue
		}

		switch m := msg.(type) {
		case *redis.Subscription:

			if m.Kind == "unsubscribe" {
				return
			}

			e.SendData(&SubscriptionResponse{
				Code: 1,
				Msg:  "Subscribed to event topic",
			})

		case *redis.Message:
			e.Send(m.Payload)
		}
	}
}

func (e *EventConsumer) Send(msg string) {
	var event struct {
		Origin          string         `json:"origin"`
		Index           uint           `json:"index"`
		Topics          pq.StringArray `json:"topics"`
		Data            string         `json:"data"`
		TransactionHash string         `json:"transaction_hash"`
		BlockHash       string         `json:"block_hash"`
		BlockNumber     uint64         `json:"block_number"`
		Timestamp       uint64         `json:"timestamp"`
	}

	_msg := []byte(data)

	if err := json.Unmarshal(msg, &event); err != nil {
		logger.S().Errorf("Failed to Decode Published event to JSON: %v", err.Error())
		return
	}

	data := make([]byte, 0)
	var err error

	if len(event.Data) != 0 {
		data, err = hex.DecodeString(event.Data[2:])
	}

	if err != nil {
		logger.S().Errorf("Failed to Decode Published event data from hex: %v", err.Error())
		return
	}

	_event := &d.Event{
		Origin:          event.Origin,
		Index:           event.Index,
		Topics:          event.Topics,
		Data:            data,
		TransactionHash: event.TransactionHash,
		BlockHash:       event.BlockHash,
		BlockNumber:     event.BlockNumber,
		Timestamp:       event.Timestamp,
	}

	var req *SubscriptionRequest

	e.TopicLock.RLock()

	for _, r := range e.Requests {
		if r.DoesMatchWithPublishedEvent(_event) {
			req = r
			break
		}
	}

	e.TopicLock.RUnlock()

	if req == nil {
		return
	}

	e.SendData(&_event)
}

func (e *EventConsumer) SendData(data interface{}) bool {
	e.ConnLock.Lock()
	defer e.ConnLock.Unlock()

	if err := e.Connection.WriteJSON(data); err != nil {
		logger.S().Errorf("Failed to send event data over client: %v", err.Error())
		return false
	}
	return true
}

func (e *EventConsumer) Unsubscribe() {
	if e.Pubsub == nil {
		logger.S().Warn("Pubsub is nil while unsubscribing from event topic")
		return
	}

	if err := e.Pubsub.Unsubscribe(context.Background(), "event"); err != nil {
		logger.S().Errorf("Failed to unsubscribe from event topic: %v", err.Error())
		return
	}

	resp := &SubscriptionResponse{
		Code: 1,
		Msg:  "Unsubscribed from event topic",
	}

	e.ConnLock.Lock()
	defer e.ConnLock.Unlock()

	if err := e.Connection.WriteJSON(resp); err != nil {
		logger.S().Errorf("Failed to send unsubscription response over client: %v", err.Error())
		return
	}
}
