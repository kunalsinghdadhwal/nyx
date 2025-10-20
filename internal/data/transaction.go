package data

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kunalsinghdadhwal/nyx/pkg/logger"
)

type Transaction struct {
	Hash            string `json:"hash" gorm:"primaryKey;column:hash"`
	From            string `json:"from" gorm:"column:from"`
	To              string `json:"to" gorm:"column:to"`
	ContractAddress string `json:"contract_address" gorm:"column:contract_address"`
	Value           string `json:"value" gorm:"column:value"`
	Data            []byte `json:"data" gorm:"column:data"`
	Gas             uint64 `json:"gas" gorm:"column:gas"`
	GasPrice        string `json:"gas_price" gorm:"column:gas_price"`
	Cost            string `json:"cost" gorm:"column:cost"`
	Nonce           uint64 `json:"nonce" gorm:"column:nonce"`
	State           uint64 `json:"state" gorm:"column:state"`
	BlockHash       string `json:"block_hash" gorm:"column:block_hash"`
	BlockNumber     uint64 `json:"block_number" gorm:"column:block_number"`
	Timestamp       uint64 `json:"timestamp" gorm:"column:timestamp"`
}

type Transactions struct {
	Transactions []*Transaction `json:"transactions"`
}

func (t *Transaction) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	data := ""

	if h := hex.EncodeToString(t.Data); h != "" {
		data = fmt.Sprintf("0x%s", h)
	}

	if !strings.HasPrefix(t.ContractAddress, "0x") {
		return []byte(fmt.Sprintf(`{"hash":%q,"from":%q,"to":%q,"value":%q,"data":%q,"gas":%d,"gasPrice":%q,"cost":%q,"nonce":%d,"state":%d,"blockHash":%q,"blockNumber":%d,"timestamp":%d}"`, t.Hash, t.From, t.To, t.Value, data, t.Gas, t.GasPrice, t.Cost, t.Nonce, t.State, t.BlockHash, t.BlockNumber, t.Timestamp)), nil
	}

	return []byte(fmt.Sprintf(
		`{hash":%q,"from":%q,"contract_address":%q,"to":%q,"value":%q,"data":%q,"gas":%d,"gasPrice":%q,"cost":%q,"nonce":%d,"state":%d,"blockHash":%q,"blockNumber":%d,"timestamp":%d}"`,
		t.Hash, t.From, t.ContractAddress, t.To, t.Value, data, t.Gas, t.GasPrice, t.Cost, t.Nonce, t.State, t.BlockHash, t.BlockNumber, t.Timestamp)), nil

}

func (t *Transaction) ToJSON() []byte {
	data, err := json.Marshal(t)

	if err != nil {
		logger.S().Errorf("Error marshaling transaction to json: %v", err.Error())
		return nil
	}

	return data
}

func (ts *Transactions) ToJSON() []byte {
	data, err := json.Marshal(ts)

	if err != nil {
		logger.S().Errorf("Error marshaling transactions to json: %v", err.Error())
		return nil
	}

	return data
}
