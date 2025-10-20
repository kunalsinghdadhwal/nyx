package data

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type SyncState struct {
	Done               uint64
	StartedAt          time.Time
	BlockCountAtStart  uint64
	MaxBlockNumAtStart uint64
	NewBlocksInserted  uint64
	LatestBlockNum     uint64
}

type StatusHolder struct {
	State *SyncState
	Mutex *sync.RWMutex
}

type RedisInfo struct {
	Client                                               *redis.Client
	BlockPublishTopic, TxPublishTopic, EventPublishTopic string
}

type ResultStatus struct {
	Success uint64
	Failure uint64
}

type Job struct {
	Client *ethclient.Client
	DB     *gorm.DB
	Redis  *RedisInfo
	Block  uint64
	Status *StatusHolder
}

type BlockChainNodeConn struct {
	RPC       *ethclient.Client
	WebSocket *ethclient.Client
}

func (s *SyncState) BlockCountInDB() uint64 {
	return s.BlockCountAtStart + s.NewBlocksInserted
}

func (s *StatusHolder) MaxBlockNumAtStart() uint64 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.State.MaxBlockNumAtStart
}

func (s *StatusHolder) SetStartedAt() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.State.StartedAt = time.Now().UTC()
}

func (s *StatusHolder) IncrementBlocksInserted() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.State.NewBlocksInserted++
}

func (s *StatusHolder) IncrementBlocksProcessed() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.State.Done++
}

func (s *StatusHolder) BlockCountInDB() uint64 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.State.BlockCountInDB()
}

func (s *StatusHolder) ElapsedTime() time.Duration {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return time.Since(s.State.StartedAt)
}

func (s *StatusHolder) Done() uint64 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.State.Done
}

func (s *StatusHolder) GetLatestBlockNum() uint64 {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.State.LatestBlockNum
}

func (s *StatusHolder) SetLatestBlockNum(num uint64) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.State.LatestBlockNum = num
}

func (r ResultStatus) Total() uint64 {
	return r.Success + r.Failure
}
