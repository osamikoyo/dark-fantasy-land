package casher

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.uber.org/zap"
)

func (c *Casher) AddWallpaperToCash(ctx context.Context, wallpaper *entity.Wallpaper) error {
	if wallpaper == nil {
		return NIL_INPUT_ERROR
	}

	c.logger.Debug("adding wallpaper to cash", zap.Any("wallpaper", wallpaper))

	key := newWallpaperKey(wallpaper.ImageName, wallpaper.Topic)

	_, err := c.client.HSet(ctx, key, wallpaper).Result()
	if err != nil {
		c.logger.Error("failed add wallpaper to cash",
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (c *Casher) GetWallpaperFromCash(ctx context.Context, imageName, topic string) (*entity.Wallpaper, error) {
	if imageName == "" || topic == "" {
		return nil, NIL_INPUT_ERROR
	}

	key := newWallpaperKey(imageName, topic)

	c.logger.Debug("fetching wallpaper", zap.String("key", key))

	res, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed get wallpaper from cash",
			zap.String("key", key),
			zap.Error(err))
		return nil, err
	}

	var wallpaper entity.Wallpaper

	if err = mapstructure.Decode(res, &wallpaper); err != nil {
		c.logger.Error("failed decode wallpaper from cash",
			zap.Any("result", res),
			zap.Error(err))
		return nil, err
	}

	return &wallpaper, nil
}

func (c *Casher) UpdateWallpaperInCash(ctx context.Context, imageName, topic, key string, value interface{}) error {
	if imageName == "" || topic == "" || key == "" {
		return NIL_INPUT_ERROR
	}

	redisKey := newWallpaperKey(imageName, topic)

	_, err := c.client.HSet(ctx, redisKey, key, value).Result()
	if err != nil {
		c.logger.Error("failed update wallpaper in hash",
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	return nil
}

func (c *Casher) DeleteWallpaperFromCash(ctx context.Context, imageName, topic string) error {
	if imageName == "" || topic == "" {
		return NIL_INPUT_ERROR
	}
	key := newWallpaperKey(imageName, topic)
	_, err := c.client.Del(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed delete wallpaper from cash",
			zap.String("key", key),
			zap.Error(err))
		return err
	}
	return nil
}
