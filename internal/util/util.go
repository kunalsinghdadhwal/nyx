package util

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"regexp"
	"strconv"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func PublicKetBytesToAddress(publickey []byte) common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publickey[1:]) // remove EC prefix 0x04
	buf = hash.Sum(nil)

	address := buf[12:]
	return common.HexToAddress(hex.EncodeToString(address))
}

func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	switch v := iaddress.(type) {
	case common.Address:
		return re.MatchString(v.Hex())
	case string:
		return re.MatchString(v)
	default:
		return false
	}
}

func IsZeroAddress(iaddress interface{}) bool {
	var address common.Address

	switch v := iaddress.(type) {
	case common.Address:
		address = v
	case string:
		address = common.HexToAddress(v)
	default:
		return false
	}

	return reflect.DeepEqual(address.Bytes(), common.FromHex("0x0000000000000000000000000000000000000000"))
}

func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)

	switch v := ivalue.(type) {
	case *big.Int:
		value = v
	case string:
		value.SetString(v, 10)
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromInt(int64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	return num.Div(mul)
}

func ToWei(iamount interface{}, decimals int) *big.Int {
	amt := decimal.NewFromInt(0)

	switch v := iamount.(type) {
	case string:
		amt, _ = decimal.NewFromString(v)
	case float64:
		amt = decimal.NewFromFloat(v)
	case int64:
		amt = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amt = v
	case *decimal.Decimal:
		amt = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	res := amt.Mul(mul)

	wei := new(big.Int)
	wei.SetString(res.String(), 10)
	return wei
}

func CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	gasLimitBig := big.NewInt(int64(gasLimit))
	return gasLimitBig.Mul(gasLimitBig, gasPrice)
}

func SigRSV(isig interface{}) (r [32]byte, s [32]byte, v uint8) {
	var sig []byte

	switch v := isig.(type) {
	case []byte:
		sig = v
	case string:
		sig, _ = hex.DecodeString(v)
	}

	sigstr := common.Bytes2Hex(sig)
	rS := sigstr[0:64]
	sS := sigstr[64:128]

	R := [32]byte{}
	S := [32]byte{}

	copy(R[:], common.FromHex(rS))
	copy(S[:], common.FromHex(sS))

	vStr := sigstr[128:130]

	vI, _ := strconv.Atoi(vStr)
	V := uint8(vI + 27)

	return R, S, V
}
func TransactionSender(block *types.Block, tx *types.Transaction) (common.Address, error) {
	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		signer = types.NewEIP2930Signer(tx.ChainId())
		sender, err = types.Sender(signer, tx)
		if err != nil {
			signer = types.NewEIP155Signer(tx.ChainId())
			sender, err = types.Sender(signer, tx)
			if err != nil {
				signer = types.HomesteadSigner{}
				sender, err = types.Sender(signer, tx)
				if err != nil {
					return common.Address{}, err
				}
			}
		}
	}
	return sender, nil
}
