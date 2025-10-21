package pubsub

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type BlockConsumer struct {
	Client     *redis.Client
	Requests   map[string]*SubscriptionRequest
	Connection *websocket.Conn
	Pubsub     *redis.PubSub
	DB         *gorm.DB
	ConnLock   *sync.Mutex
	TopicLock  *sync.RWMutex
}

func (b *BlockConsumer) Subscribe() {
	b.Pubsub = b.Client.Subscribe(context.Background(), "block")
}
