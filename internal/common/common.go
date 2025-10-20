package common

import (
	"errors"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

func StringifyEventTopics(data []common.Hash) []string {
	buffer := make([]string, len(data))

	for i := 0; i < len(data); i++ {
		buffer[i] = data[i].Hex()
	}

	return buffer
}

func CreateEventTopicMap(topics []string) map[uint8]string {
	topicMap := make(map[uint8]string)

	if topics[0] != "" {
		topicMap[0] = topics[0]
	}

	if topics[1] != "" {
		topicMap[1] = topics[1]
	}

	if topics[2] != "" {
		topicMap[2] = topics[2]
	}

	if topics[3] != "" {
		topicMap[3] = topics[3]
	}

	return topicMap
}

func RangeChecker(from string, to string, limit uint64) (uint64, uint64, error) {
	fromInt, err := strconv.Atoi(from)

	if err != nil {
		return 0, 0, errors.New("[Range Checker] Failed to parse 'from' parameter")
	}

	toInt, err := strconv.Atoi(to)

	if err != nil {
		return 0, 0, errors.New("[Range Checker] Failed to parse 'to' parameter")
	}

	if uint64(toInt-fromInt) > limit {
		return 0, 0, errors.New("[Range Checker] Range exceeds maximum limit")
	}

	return uint64(fromInt), uint64(toInt), nil
}
