package casher

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.uber.org/zap"
)

func (c *Casher) AddMemToCash(ctx context.Context, mem *entity.Mem) error {
	if mem == nil {
		return NIL_INPUT_ERROR
	}

	c.logger.Debug("adding mem to cash", zap.Any("mem", mem))

	key := newMemKey(mem.ImageName, mem.Author)

	_, err := c.client.HSet(ctx, key, mem).Result()
	if err != nil {
		c.logger.Error("failed add mem to cash",
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (c *Casher) GetMemFromCash(ctx context.Context, imageName, author string) (*entity.Mem, error) {
	if imageName == "" || author == "" {
		return nil, NIL_INPUT_ERROR
	}

	key := newMemKey(imageName, author)

	c.logger.Debug("fetching mem", zap.String("key", key))

	res, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed get mem from cash",
			zap.String("key", key),
			zap.Error(err))
		return nil, err
	}

	var mem entity.Mem

	if err = mapstructure.Decode(res, &mem); err != nil {
		c.logger.Error("failed decode mem from cash",
			zap.Any("result", res),
			zap.Error(err))
		return nil, err
	}

	return &mem, nil
}

func (c *Casher) UpdateMemInCash(ctx context.Context, imageName, author, key string, value interface{}) error {
	if imageName == "" || author == "" || key == "" {
		return NIL_INPUT_ERROR
	}

	redisKey := newMemKey(imageName, author)

	_, err := c.client.HSet(ctx, redisKey, key, value).Result()
	if err != nil {
		c.logger.Error("failed update mem in hash",
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (c *Casher) DeleteMemFromCash(ctx context.Context, imageName, author string) error {
	if imageName == "" || author == "" {
		return NIL_INPUT_ERROR
	}
	key := newMemKey(imageName, author)
	_, err := c.client.Del(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed delete mem from cash",
			zap.String("key", key),
			zap.Error(err))
		return err
	}
	return nil
}
