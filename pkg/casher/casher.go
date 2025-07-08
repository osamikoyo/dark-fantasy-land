package casher

import (
	"errors"

	"github.com/osamikoyo/dark-fantasy-land/pkg/logger"
	"github.com/redis/go-redis/v9"
)

var NIL_INPUT_ERROR = errors.New("input value in nil")

type Casher struct {
	client *redis.Client
	logger *logger.Logger
}

func NewCasher(client *redis.Client, logger *logger.Logger) *Casher {
	return &Casher{
		client: client,
		logger: logger,
	}
}
