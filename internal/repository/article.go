package repository

import (
	"context"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *Repository) UpdateArticle(ctx context.Context, filter, update map[string]interface{}) error {
	r.logger.Debug("updating article",
		zap.Any("filter", filter),
		zap.Any("update", update))

	res, err := r.articlesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("failed update article", zap.Error(err))

		return err
	}

	r.logger.Info("article updated",
		zap.Int64("modifed_count", res.ModifiedCount))

	return nil
}

func (r *Repository) DeleteArticle(ctx context.Context, filter map[string]string) error {
	r.logger.Debug("deleting article",
		zap.Any("filter", filter))

	_, err := r.articlesColl.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("failed delete article",
			zap.Error(err))

		return err
	}

	r.logger.Info("deleted article",
		zap.Any("filter", filter))

	return nil
}

func (r *Repository) GetArticle(ctx context.Context, filter map[string]interface{}) ([]entity.Article, error) {
	r.logger.Debug("fetching article",
		zap.Any("filter", filter))

	res, err := r.articlesColl.Find(ctx, filter)
	if err != nil {
		r.logger.Error("failed fetch articles",
			zap.Any("filter", filter),
			zap.Error(err))

		return nil, err
	}

	var articles []entity.Article

	for res.Next(ctx) {
		var article entity.Article

		if err = res.Decode(&article); err != nil {
			r.logger.Warn("failed decode article", zap.Error(err))

			continue
		}

		articles = append(articles, article)
	}

	if err = res.Err(); err != nil {
		r.logger.Error("error from fetch response", zap.Error(err))

		return nil, err
	}

	return articles, nil
}

func (r *Repository) GetArticlesLimited(ctx context.Context, filter map[string]interface{}, limit int64) ([]entity.Article, error) {
	r.logger.Debug("fetching limited articles", zap.Any("filter", filter), zap.Int64("limit", limit))

	findOptions := options.Find()
	findOptions.SetLimit(limit)

	res, err := r.articlesColl.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("failed fetch limited articles", zap.Any("filter", filter), zap.Error(err))
		return nil, err
	}

	var articles []entity.Article
	for res.Next(ctx) {
		var article entity.Article
		if err = res.Decode(&article); err != nil {
			r.logger.Warn("failed decode article", zap.Error(err))
			continue
		}
		articles = append(articles, article)
	}

	if err = res.Err(); err != nil {
		r.logger.Error("error from fetch response", zap.Error(err))
		return nil, err
	}

	return articles, nil
}
