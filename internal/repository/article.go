package repository

import (
	"context"
	"fmt"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) CreateArticle(ctx context.Context, article *entity.Article) error {
	r.logger.Debug("creating article", zap.Any("article", article))

	res, err := r.articlesColl.InsertOne(ctx, article)
	if err != nil {
		r.logger.Error("failed create article", zap.String("title", article.Title), zap.Error(err))
		return fmt.Errorf("create article: %w", ErrInsertFailed)
	}

	r.logger.Info("article created", zap.String("inserted_id", fmt.Sprintf("%v", res.InsertedID)))
	return nil
}

func (r *Repository) UpdateArticle(ctx context.Context, filter, update map[string]interface{}) error {
	r.logger.Debug("updating article", zap.Any("filter", filter), zap.Any("update", update))

	res, err := r.articlesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("failed update article", zap.Error(err))
		return fmt.Errorf("update article: %w", ErrUpdateFailed)
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	r.logger.Info("article updated", zap.Int64("modifed_count", res.ModifiedCount))
	return nil
}

func (r *Repository) DeleteArticle(ctx context.Context, filter map[string]interface{}) error {
	r.logger.Debug("deleting article", zap.Any("filter", filter))
	res, err := r.articlesColl.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("failed delete article", zap.Error(err))
		return fmt.Errorf("delete article: %w", ErrDeleteFailed)
	}

	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	r.logger.Info("deleted article", zap.Any("filter", filter))
	return nil
}

func (r *Repository) GetArticle(ctx context.Context, filter map[string]interface{}) (*entity.Article, error) {
	r.logger.Debug("fetching single article", zap.Any("filter", filter))

	res := r.articlesColl.FindOne(ctx, filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			r.logger.Warn("article not found", zap.Any("filter", filter))
			return nil, ErrNotFound
		}
		r.logger.Error("failed to get article", zap.Error(res.Err()))
		return nil, fmt.Errorf("get article: %w", res.Err())
	}

	var article entity.Article
	if err := res.Decode(&article); err != nil {
		r.logger.Warn("failed decode article", zap.Error(err))
		return nil, fmt.Errorf("decode article: %w", ErrDecodeFailed)
	}

	r.logger.Info("article fetched", zap.Any("article", article))
	return &article, nil
}

func (r *Repository) GetArticlesLimited(ctx context.Context, filter map[string]interface{}, limit int64) ([]entity.Article, error) {
	r.logger.Debug("fetching limited articles", zap.Any("filter", filter), zap.Int64("limit", limit))

	findOptions := options.Find()
	findOptions.SetLimit(limit)

	res, err := r.articlesColl.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("failed fetch limited articles", zap.Any("filter", filter), zap.Error(err))
		return nil, fmt.Errorf("get limited articles: %w", ErrNotFound)
	}

	var articles []entity.Article
	for res.Next(ctx) {
		var article entity.Article
		if err = res.Decode(&article); err != nil {
			r.logger.Warn("failed decode article", zap.Error(err))
			return nil, fmt.Errorf("decode article: %w", ErrDecodeFailed)
		}
		articles = append(articles, article)
	}

	if err = res.Err(); err != nil {
		r.logger.Error("error from fetch response", zap.Error(err))
		return nil, fmt.Errorf("parse articles: %w", ErrDecodeFailed)
	}

	if len(articles) == 0 {
		return nil, ErrNoDocuments
	}

	return articles, nil
}
