package block

import (
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kunalsinghdadhwal/nyx/pkg/logger"
	"gorm.io/gorm"
)

func ProcessBlock(client *ethclient.Client, block *types.Block, db *gorm.DB, redis )