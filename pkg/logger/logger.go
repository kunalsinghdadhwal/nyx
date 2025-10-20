package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	once sync.Once
)

func Init(env string) {
	once.Do(func() {
		var err error

		if env == "prod" {
			log, err = zap.NewProduction()
		} else {
			log, err = zap.NewDevelopment()
		}

		if err != nil {
			panic("Failed to initialize logger: " + err.Error())
		}
	})
}

func L() *zap.Logger {
	if log == nil {
		panic("Logger not initialized. Call Init() first.")
	}
	return log
}

func S() *zap.SugaredLogger {
	return L().Sugar()
}
