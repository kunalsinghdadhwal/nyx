package client

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"github.com/kunalsinghdadhwal/nyx/pkg/logger"
)

func getClient(isRPC bool) *ethclient.Client {
	log := logger.S()

	var client *ethclient.Client
	var err error

	var url string
	if isRPC {
		url = os.Getenv("RPC_URL")
	} else {
		url = os.Getenv("WS_URL")
	}

	client, err = ethclient.Dial(url)

	if err != nil {
		log.Fatalf("Failed to connect to blockchain: %s\n", err.Error())
	}

	return client

}

func getRedisClient() *redis.Client {
	var options *redis.Options

	if os.Getenv("REDIS_PASSWORD") != "" {
		options = &redis.Options{
			Network:  os.Getenv("REDIS_CONN"),
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		}
	} else {
		options = &redis.Options{
			Network: os.Getenv("REDIS_CONN"),
			Addr:    os.Getenv("REDIS_ADDR"),
			DB:      0,
		}
	}

	client := redis.NewClient(options)

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.S().Fatalf("Failed to connect to Redis: %s\n", err.Error())
		return nil
	}
	return client
}
