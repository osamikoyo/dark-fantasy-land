package casher

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.uber.org/zap"
)

func (c *Casher) AddArticleToCash(ctx context.Context, article *entity.Article) error {
	if article != nil {
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
	if author == "" && title == "" {
		return nil, NIL_INPUT_ERROR
	}

	key := newArticleKey(author, title)

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
