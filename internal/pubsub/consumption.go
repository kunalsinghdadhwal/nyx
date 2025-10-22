package pubsub

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Consumer interface {
	Subscribe()
	Listen()
	Send(data string)
	SendData(data interface{}) bool
	Unsubscribe()
}

func NewBlockConsumer(client *redis.Client, requests map[string]*SubscriptionRequest, conn *websocket.Conn, db *gorm.DB, connLock *sync.Mutex, topicLock *sync.RWMutex) *BlockConsumer {
	consumer := BlockConsumer{
		Client:     client,
		Requests:   requests,
		Connection: conn,
		DB:         db,
		ConnLock:   connLock,
		TopicLock:  topicLock,
	}

	consumer.Subscribe()
	go consumer.Listen()

	return &consumer
}

func NewTransactionConsumer(client *redis.Client, requests map[string]*SubscriptionRequest, conn *websocket.Conn, db *gorm.DB, connLock *sync.Mutex, topicLock *sync.RWMutex) *TransactionConsumer {
	consumer := TransactionConsumer{
		Client:     client,
		Requests:   requests,
		Connection: conn,
		DB:         db,
		ConnLock:   connLock,
		TopicLock:  topicLock,
	}

	consumer.Subscribe()
	go consumer.Listen()

	return &consumer
}

func NewEventConsumer(client *redis.Client, requests map[string]*SubscriptionRequest, conn *websocket.Conn, db *gorm.DB, connLock *sync.Mutex, topicLock *sync.RWMutex) *EventConsumer {
	consumer := EventConsumer{
		Client:     client,
		Requests:   requests,
		Connection: conn,
		DB:         db,
		ConnLock:   connLock,
		TopicLock:  topicLock,
	}

	consumer.Subscribe()
	go consumer.Listen()

	return &consumer
}
