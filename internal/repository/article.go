package repository

import (
	"context"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.uber.org/zap"
)

func (r *Repository) CreateArticle(ctx context.Context, article *entity.Article) error {
	r.logger.Debug("creating article", zap.Any("article", article))

	res, err := r.articlesColl.InsertOne(ctx, article)
	if err != nil {
		r.logger.Error("failed create article",
			zap.String("title", article.Title),
			zap.Error(err))

		return err
	}

	r.logger.Info("article created",
		zap.String("inserted_id", res.InsertedID.(string)))

	return nil
}
