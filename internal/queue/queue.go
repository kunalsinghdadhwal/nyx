package queue

import (
	"context"
	"math"
	"os"
	"strconv"
	"time"
)

type Block struct {
	UnconfirmedProgress bool
	Published           bool
	UnconfirmedDone     bool
	ConfirmedProgress   bool
	ConfirmedDone       bool
	LastAttempted       time.Time
	Delay               time.Duration
}

type Request struct {
	BlockNumber  uint64
	ResponseChan chan bool
}

type Update Request

type Next struct {
	ResponseChan chan struct {
		Status bool
		Number uint64
	}
}

type StatResponse struct {
	UnconfirmedProgress uint64
	UnconfirmedWaiting  uint64
	ConfirmedProgress   uint64
	ConfirmedWaiting    uint64
	Total               uint64
}

type Stat struct {
	ResponseChan chan StatResponse
}

type BlockProcessorQueue struct {
	Blocks               map[uint64]*Block
	StartedWith          uint64
	TotalInserted        uint64
	LatestBlock          uint64
	Total                uint64
	PutChan              chan Request
	CanPublishChan       chan Request
	PublishedChan        chan Request
	InsertedChan         chan Request
	UnconfimedFailedChan chan Request
	UnconfirmedDoneChan  chan Request
	ConfirmedFailedChan  chan Request
	StatChan             chan Stat
	LatestChan           chan Update
	UnconfirmedNextChan  chan Next
	ConfirmedNextChan    chan Next
}

func (b *Block) SetDelay() {
	b.Delay = time.Duration(int64(math.Round(b.Delay.Seconds()*(1.0+math.Sqrt(5.0))/2))%3600) * time.Second
}

func (b *Block) ResetDelay() {
	b.Delay = time.Duration(1) * time.Second
}

func (b *Block) SetLastAttempted() {
	b.LastAttempted = time.Now().UTC()
}

func (b *Block) CanAttempt() bool {
	return time.Now().UTC().After(b.LastAttempted.Add(b.Delay))
}

func New(startedWith uint64) *BlockProcessorQueue {
	return &BlockProcessorQueue{
		Blocks:               make(map[uint64]*Block),
		StartedWith:          startedWith,
		TotalInserted:        0,
		LatestBlock:          0,
		Total:                0,
		PutChan:              make(chan Request, 128),
		CanPublishChan:       make(chan Request, 128),
		PublishedChan:        make(chan Request, 128),
		InsertedChan:         make(chan Request, 128),
		UnconfimedFailedChan: make(chan Request, 128),
		UnconfirmedDoneChan:  make(chan Request, 128),
		ConfirmedFailedChan:  make(chan Request, 128),
		StatChan:             make(chan Stat, 1),
		LatestChan:           make(chan Update, 1),
		UnconfirmedNextChan:  make(chan Next, 1),
		ConfirmedNextChan:    make(chan Next, 1),
	}
}


func (q *BlockProcessorQueue) Put(block uint64) bool {
	resp := make(chan bool)

	req := Request{
		BlockNumber:  block,
		ResponseChan: resp,
	}

	q.PutChan <- req

	return <-resp
}

func (q *BlockProcessorQueue) CanPublish(block uint64) bool {
	resp := make(chan bool)

	req := Request{
		BlockNumber:  block,
		ResponseChan: resp,
	}
	
	q.CanPublishChan <- req
	return <-resp
}

func (q *BlockProcessorQueue) Published(block uint64) bool {
	resp := make(chan bool)

	req := Request{
		BlockNumber:  block,
		ResponseChan: resp,
	}

	q.PublishedChan <- req

	return <-resp
}

func (q *BlockProcessorQueue) Inserted(block uint64) bool {
	resp := make(chan bool)

	req := Request{
		BlockNumber:  block,
		ResponseChan: resp,
	}

	q.InsertedChan <- req

	return <-resp
}

func (q *BlockProcessorQueue) UnconfirmedFailed(block uint64) bool {
	resp := make(chan bool)

	req := Request{
		BlockNumber:  block,
		ResponseChan: resp,
	}

	q.UnconfimedFailedChan <- req

	return <-resp
}

func (q *BlockProcessorQueue) UnconfirmedDone(block uint64) bool {
	resp := make(chan bool)

	req := Request{
		BlockNumber:  block,
		ResponseChan: resp,
	}

	q.UnconfirmedDoneChan <- req

	return <-resp
}

func (q *BlockProcessorQueue) ConfirmedFailed(block uint64) bool {
	resp := make(chan bool)

	req := Request{
		BlockNumber:  block,
		ResponseChan: resp,
	}

	q.ConfirmedFailedChan <- req

	return <-resp
}

func (q *BlockProcessorQueue) ConfirmedDone(block uint64) bool {
	resp := make(chan bool)

	req := Request{
		BlockNumber:  block,
		ResponseChan: resp,
	}

	q.ConfirmedFailedChan <- req

	return <-resp
}

func (q *BlockProcessorQueue) Stat() StatResponse {
	resp := make(chan StatResponse)

	req := Stat{
		ResponseChan: resp,
	}

	q.StatChan <- req

	return <-resp
}

func (q *BlockProcessorQueue) Latest(num uint64) bool {
	resp := make(chan bool)

	req := Update{
		BlockNumber:  num,
		ResponseChan: resp,
	}

	q.LatestChan <- req

	return <-resp	
}

func (q *BlockProcessorQueue) UnconfirmedNext() (uint64, bool) {
	resp := make(chan struct {
		Status bool
		Number uint64
	})

	req := Next{
		ResponseChan: resp,
	}

	q.UnconfirmedNextChan <- req

	result := <-resp

	return result.Number, result.Status
}

func (q *BlockProcessorQueue) ConfirmedNext() (uint64, bool) {
	resp := make(chan struct {
		Status bool
		Number uint64
	})

	req := Next{
		ResponseChan: resp,
	}
	q.ConfirmedNextChan <- req

	result := <-resp

	return result.Number, result.Status
}

func (q *BlockProcessorQueue) CanBeConfirmed(block uint64) bool {
	var blockConfirmations int = 0
	
	if os.Getenv("BLOCK_CONFIRMATIONS") != "" {
		blockConfirmationsEnv, err := strconv.Atoi(os.Getenv("BLOCK_CONFIRMATIONS"))
		if err != nil {
			blockConfirmations = 0
		} else {
			blockConfirmations = blockConfirmationsEnv
		}
	}

	if q.LatestBlock < uint64(blockConfirmations) {
		return false
	}

	return q.LatestBlock - uint64(blockConfirmations) >= block
}

func (q *BlockProcessorQueue) TotalBlocks() uint64 {
	return q.Total
}

func (q *BlockProcessorQueue) Start(ctx context.Context) {
	for {
		select {
			case <-ctx.Done():
				return

			case req := <-q.PutChan:
				if _, ok := q.Blocks[req.BlockNumber]; ok {
					req.ResponseChan <- false
					break
				}

				q.Blocks[req.BlockNumber] = &Block{
					UnconfirmedProgress: true,
					LastAttempted: time.Now().UTC(),
					Delay: time.Duration(1) * time.Second,
				}
				req.ResponseChan <- true
			
			case req := <-q.CanPublishChan:
				block, ok := q.Blocks[req.BlockNumber]
				if !ok {
					req.ResponseChan <- false
					break
				}
				req.ResponseChan <- !block.Published
		
			
		}
	}
}
