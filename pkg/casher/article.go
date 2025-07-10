package casher

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.uber.org/zap"
)

func (c *Casher) AddArticleToCash(ctx context.Context, article *entity.Article) error {
	if article == nil {
		return NIL_INPUT_ERROR
	}

	c.logger.Debug("adding article to cash", zap.Any("article", article))

	key := newArticleKey(article.Author, article.Title)

	_, err := c.client.HSet(ctx, key, article).Result()
	if err != nil {
		c.logger.Error("failed add articel to cash",
			zap.String("key", key),
			zap.Error(err))

		return err
	}

	return nil
}

func (c *Casher) GetArticleFromCash(ctx context.Context, author, title string) (*entity.Article, error) {
	if author == "" || title == "" {
		return nil, NIL_INPUT_ERROR
	}

	key := newArticleKey(author, title)

	c.logger.Debug("fetching article", zap.String("key", key))

	res, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		c.logger.Error("failed get article from cash",
			zap.String("key", key),
			zap.Error(err))

		return nil, err
	}

	var article entity.Article

	if err = mapstructure.Decode(res, &article); err != nil {
		c.logger.Error("failed decode article from cash",
			zap.Any("result", res),
			zap.Error(err))

		return nil, err
	}

	return &article, nil
}

func (c *Casher) UpdateArticleInCash(ctx context.Context, author, title, key string, value interface{}) error {
	if author == "" || title == "" || key == "" {
		return NIL_INPUT_ERROR
	}

	redisKey := newArticleKey(author, title)

	_, err := c.client.HSet(ctx, redisKey, key, value).Result()
	if err != nil {
		c.logger.Error("failed update article in hash",
			zap.String("key", key),
			zap.Error(err))

		return err
	}

	return nil
}

func (c *Casher) DeleteArticleFromCash(ctx context.Context, author, title string) error {
	if author == "" || title == "" {
		return NIL_INPUT_ERROR
	}

	redisKey := newArticleKey(author, title)

	c.logger.Debug("deleting articles in cash", zap.String("key", redisKey))

	_, err := c.client.Del(ctx, redisKey).Result()
	if err != nil {
		c.logger.Error("failed to delete article",
			zap.String("key", redisKey),
			zap.Error(err))

		return err
	}

	return nil
}
