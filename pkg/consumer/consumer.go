package consumer

import (
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"github.com/osamikoyo/dark-fantasy-land/internal/service"
	"github.com/osamikoyo/dark-fantasy-land/pkg/logger"
	"go.uber.org/zap"
)

type Consumer struct {
	logger  *logger.Logger
	service *service.Service
	client  *nats.Conn
}

func NewConsumer(logger *logger.Logger, service *service.Service, client *nats.Conn) *Consumer {
	return &Consumer{
		logger:  logger,
		service: service,
		client:  client,
	}
}

func (c *Consumer) SubscribeAll() error {
	_, err := c.client.Subscribe("censored_articles", func(msg *nats.Msg) {
		var req entity.Request[*entity.Article]

		if err := sonic.Unmarshal(msg.Data, &req); err != nil {
			c.logger.Error("failed unmarshal message body",
				zap.String("subject", msg.Subject),
				zap.Error(err))
			return
		}

		if err := c.service.CreateArticle(req.Payload); err != nil {
			c.logger.Error("failed add article", zap.Error(err))
		}
	})
	if err != nil {
		c.logger.Error("failed subscribe on censored_articles", zap.Error(err))

		return err
	}

	_, err = c.client.Subscribe("censored_mems", func(msg *nats.Msg) {
		var req entity.Request[*entity.Mem]

		if err := sonic.Unmarshal(msg.Data, &req); err != nil {
			c.logger.Error("failed unmarshal message body",
				zap.String("subject", msg.Subject),
				zap.Error(err))

			return
		}

		if err := c.service.CreateMem(req.Payload); err != nil {
			c.logger.Error("failed create mem",
				zap.Any("mem", req.Payload),
				zap.Error(err))

			return
		}
	})
	if err != nil{
		c.logger.Error("failed subscribe on censored_mems", zap.Error(err))

		return err
	}

	return nil
}
