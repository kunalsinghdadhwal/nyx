package client

import (
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
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
