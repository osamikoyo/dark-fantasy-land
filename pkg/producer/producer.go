package producer

import (
	"errors"

	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/osamikoyo/dark-fantasy-land/pkg/logger"
	"go.uber.org/zap"
)

type Producer struct {
	client *nats.Conn
	logger *logger.Logger
}

func NewProducer(client *nats.Conn, logger *logger.Logger) *Producer {
	return &Producer{
		client: client,
		logger: logger,
	}
}

func (p *Producer) SendToCensor(queue string, value interface{}) error {
	if value == nil {
		return errors.New("nil input")
	}

	payload, err := sonic.Marshal(value)
	if err != nil {
		p.logger.Error("failed marshal value",
			zap.Any("value", value),
			zap.String("queue", queue),
			zap.Error(err))

		return err
	}

	if err = p.client.Publish(queue, payload); err != nil {
		p.logger.Error("failed publish value",
			zap.Any("value", value),
			zap.String("queue", queue),
			zap.Error(err))

		return err
	}

	return nil
}
