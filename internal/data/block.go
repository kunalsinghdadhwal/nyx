package data

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/kunalsinghdadhwal/nyx/pkg/logger"
)

type Block struct {
	Hash                string  `json:"hash" gorm:"column:hash;primaryKey"`
	Number              uint64  `json:"number" gorm:"column:number"`
	Time                uint64  `json:"time" gorm:"column:time"`
	ParentHash          string  `json:"parent_hash" gorm:"column:parent_hash"`
	Difficulty          string  `json:"difficulty" gorm:"column:difficulty"`
	GasUsed             uint64  `json:"gas_used" gorm:"column:gas_used"`
	GasLimit            uint64  `json:"gas_limit" gorm:"column:gas_limit"`
	Nonce               string  `json:"nonce" gorm:"column:nonce"`
	Miner               string  `json:"miner" gorm:"column:miner"`
	Size                float64 `json:"size" gorm:"column:size"`
	StateRootHash       string  `json:"state_root_hash" gorm:"column:state_root_hash"`
	UncleHash           string  `json:"uncle_hash" gorm:"column:uncle_hash"`
	TransactionRootHash string  `json:"transaction_root_hash" gorm:"column:transaction_root_hash"`
	ReceiptRootHash     string  `json:"receipt_root_hash" gorm:"column:receipt_root_hash"`
	ExtraData           []byte  `json:"extra_data" gorm:"column:extra_data"`
}

type Blocks struct {
	Blocks []*Block `json:"blocks"`
}

func (b *Block) MarshalBinary() ([]byte, error) {
	jsonData, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (b *Block) MarshalJSON() ([]byte, error) {
	var extraData string
	if h := hex.EncodeToString(b.ExtraData); h != "" {
		extraData = fmt.Sprintf("0x%s", h)
	}

	return []byte(fmt.Sprintf(`{"hash":%q,"number":%d,"time":%d,"parentHash":%q,"difficulty":%q,"gasUsed":%d,"gasLimit":%d,"nonce":%q,"miner":%q,"size":%f,"stateRootHash":%q,"uncleHash":%q,"txRootHash":%q,"receiptRootHash":%q,"extraData":%q}`,
		b.Hash,
		b.Number,
		b.Time,
		b.ParentHash,
		b.Difficulty,
		b.GasUsed,
		b.GasLimit,
		b.Nonce,
		b.Miner,
		b.Size,
		b.StateRootHash,
		b.UncleHash,
		b.TransactionRootHash,
		b.ReceiptRootHash,
		extraData)), nil
}

func (b *Block) ToJSON() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		logger.S().Errorf("failed to marshal block to JSON: %v", err.Error())
		return nil
	}

	return data
}

func (b *Blocks) ToJSON() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		logger.S().Errorf("failed to marshal blocks to JSON: %v", err.Error())
		return nil
	}

	return data
}
