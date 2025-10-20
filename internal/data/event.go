package data

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kunalsinghdadhwal/nyx/pkg/logger"
	"github.com/lib/pq"
)

type Event struct {
	Origin          string         `gorm:"column:origin"`
	Index           uint           `gorm:"column:index"`
	Topics          pq.StringArray `gorm:"column:topics;type:text[]"`
	Data            []byte         `gorm:"column:data"`
	TransactionHash string         `gorm:"column:transaction_hash"`
	BlockHash       string         `gorm:"column:block_hash"`
	BlockNumber     uint64         `gorm:"column:block_number"`
	Timestamp       uint64         `gorm:"column:timestamp"`
}

type Events struct {
	Events []*Event `json:"events"`
}

func (e *Event) MarshalBinary() (data []byte, err error) {
	return e.MashalJSON()
}

func (e *Event) MashalJSON() ([]byte, error) {
	data := ""

	if h := hex.EncodeToString(e.Data); h != "" && h != strings.Repeat("0", 64) {
		data = fmt.Sprintf("0x%s", h)
	}

	topics := strings.Join(strings.Fields(fmt.Sprintf("%q", e.Topics)), ",")

	return []byte(fmt.Sprintf(`{"origin":%q,"index":%d,"topics":%v,"data":%q,"txHash":%q,"blockHash":%q,"blockNumber":%d,"timestamp":%d}`,
		e.Origin,
		e.Index,
		topics,
		data,
		e.TransactionHash,
		e.BlockHash,
		e.BlockNumber,
		e.Timestamp)), nil
}

func (e *Event) ToJSON() []byte {
	data, err := json.Marshal(e)

	if err != nil {
		logger.S().Errorf("Error marshaling event to JSON: %v", err.Error())
		return nil
	}

	return data
}

func (e *Events) ToJSON() []byte {
	data, err := json.Marshal(e)

	if err != nil {
		logger.S().Errorf("Error marshaling events to JSON: %v", err.Error())
		return nil
	}

	return data
}
