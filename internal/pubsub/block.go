package pubsub

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/kunalsinghdadhwal/nyx/pkg/logger"
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

func (b *BlockConsumer) Listen() {
	for {
		msg, err := b.Pubsub.ReceiveTimeout(context.Background(), time.Duration(1)*time.Second)
		if err != nil {
			continue
		}

		switch m := msg.(type) {
		case *redis.Subscription:

			if m.Kind == "unsubscribe" {
				return
			}

			b.SendData(&SubscriptionResponse{
				Code: 1,
				Msg:  "Subscribed to block topic",
			})

		case *redis.Message:
			b.Send(m.Payload)
		}
	}
}

func (b *BlockConsumer) Send(data string) {
	var req *SubscriptionRequest

	b.TopicLock.RLock()

	for _, r := range b.Requests {
		req = r
		break
	}

	b.TopicLock.RUnlock()

	if req == nil {
		return
	}

	var block struct {
		Hash                string  `json:"hash"`
		Number              uint64  `json:"number"`
		Time                uint64  `json:"time" `
		ParentHash          string  `json:"parent_hash"`
		Difficulty          string  `json:"difficulty"`
		GasUsed             uint64  `json:"gas_used" `
		GasLimit            uint64  `json:"gas_limit"`
		Nonce               string  `json:"nonce"`
		Miner               string  `json:"miner"`
		Size                float64 `json:"size"`
		StateRootHash       string  `json:"state_root_hash"`
		UncleHash           string  `json:"uncle_hash"`
		TransactionRootHash string  `json:"transaction_root_hash"`
		ReceiptRootHash     string  `json:"receipt_root_hash"`
		ExtraData           string  `json:"extra_data"`
	}

	msg := []byte(data)

	err := json.Unmarshal(msg, &block)
	if err != nil {
		logger.S().Errorf("Failed to Decode Published block to JSON: %v", err.Error())
		return
	}

	b.SendData(&block)
}

func (b *BlockConsumer) SendData(data interface{}) bool {
	b.ConnLock.Lock()
	defer b.ConnLock.Unlock()

	if err := b.Connection.WriteJSON(data); err != nil {
		logger.S().Errorf("Failed to send block data over client: %v", err.Error())
		return false
	}
	return true
}

func (b *BlockConsumer) Unsubscribe() {
	if b.Pubsub == nil {
		logger.S().Warn("Pubsub is nil, cannot unsubscribe")
		return
	}

	if err := b.Pubsub.Unsubscribe(context.Background(), "block"); err != nil {
		logger.S().Errorf("Failed to unsubscribe from block topic: %v", err.Error())
		return
	}

	resp := &SubscriptionResponse{
		Code: 1,
		Msg:  "Unsubscribed from block topic",
	}

	b.ConnLock.Lock()
	defer b.ConnLock.Unlock()

	if err := b.Connection.WriteJSON(resp); err != nil {
		logger.S().Errorf("Failed to send unsubscribe confirmation over client: %v", err.Error())
		return
	}
}
