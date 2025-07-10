package casher

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.uber.org/zap"
)

func (c *Casher) AddNewToCash(ctx context.Context, n *entity.New) error {
	if n == nil {
		return NIL_INPUT_ERROR
	}

	c.logger.Debug("adding new to cash", zap.Any("new", n))

	key := newNewKey(n.Title, n.Author)

	_, err := c.client.HSet(ctx, key, n).Result()
	if err != nil {
		c.logger.Error("failed add new to cash",
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (c *Casher) GetNewFromCash(ctx context.Context, title, author string) (*entity.New, error) {
	if title == "" || author == "" {
		return nil, NIL_INPUT_ERROR
	}

	key := newNewKey(title, author)

	c.logger.Debug("fetching new", zap.String("key", key))

	res, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed get new from cash",
			zap.String("key", key),
			zap.Error(err))
		return nil, err
	}

	var n entity.New

	if err = mapstructure.Decode(res, &n); err != nil {
		c.logger.Error("failed decode new from cash",
			zap.Any("result", res),
			zap.Error(err))
		return nil, err
	}

	return &n, nil
}

func (c *Casher) UpdateNewInCash(ctx context.Context, title, author, key string, value interface{}) error {
	if title == "" || author == "" || key == "" {
		return NIL_INPUT_ERROR
	}

	redisKey := newNewKey(title, author)

	_, err := c.client.HSet(ctx, redisKey, key, value).Result()
	if err != nil {
		c.logger.Error("failed update new in hash",
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (c *Casher) DeleteNewFromCash(ctx context.Context, title, author string) error {
	if title == "" || author == "" {
		return NIL_INPUT_ERROR
	}
	key := newNewKey(title, author)
	_, err := c.client.Del(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed delete new from cash",
			zap.String("key", key),
			zap.Error(err))
		return err
	}
	return nil
}
