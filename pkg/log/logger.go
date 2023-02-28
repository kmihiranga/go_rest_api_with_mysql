package log

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var logger *zap.Logger
var once sync.Once

func GetLogger() *zap.Logger {
	once.Do(func() {
		content, err := os.ReadFile("./ops/configs/logger.yaml")
		if err != nil {
			panic(fmt.Sprintf("error: %v", err))
		}

		var cfg zap.Config

		err = yaml.Unmarshal(content, &cfg)
		if err != nil {
			panic(fmt.Sprintf("error: %v", err))
		}

		logger, err = cfg.Build()
		if err != nil {
			panic(err)
		}
		defer logger.Sync()

		logger.Info("Logger initialized.")
	})
	return logger
}
