package repository

import (
	"context"
	"fmt"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) CreateNew(ctx context.Context, New *entity.New) error {
	res, err := r.newsColl.InsertOne(ctx, New)
	if err != nil {
		r.logger.Error("failed create new", zap.Error(err))
		return fmt.Errorf("create new: %w", ErrInsertFailed)
	}

	r.logger.Info("new created",
		zap.String("title", New.Title),
		zap.String("inserted_id", fmt.Sprintf("%v", res.InsertedID)))

	return nil
}

func (r *Repository) UpdateNew(ctx context.Context, filter, update map[string]interface{}) error {
	r.logger.Debug("updating new", zap.Any("filter", filter), zap.Any("update", update))

	res, err := r.newsColl.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("failed update one", zap.Error(err))
		return fmt.Errorf("update new: %w", ErrUpdateFailed)
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	r.logger.Info("update one new", zap.Int64("matched_count", res.MatchedCount))

	return nil
}

func (r *Repository) GetNews(ctx context.Context, filter map[string]interface{}) (*entity.New, error) {
    r.logger.Debug("fetching single news", zap.Any("filter", filter))

    res := r.newsColl.FindOne(ctx, filter)
    if res.Err() != nil {
        if res.Err() == mongo.ErrNoDocuments {
            r.logger.Warn("news not found", zap.Any("filter", filter))
            return nil, ErrNotFound
        }
        r.logger.Error("failed to get news", zap.Error(res.Err()))
        return nil, fmt.Errorf("get news: %w", res.Err())
    }

    var n entity.New
    if err := res.Decode(&n); err != nil {
        r.logger.Warn("failed decode news", zap.Error(err))
        return nil, fmt.Errorf("decode news: %w", ErrDecodeFailed)
    }

    r.logger.Info("news fetched", zap.Any("news", n))
    return &n, nil
}

func (r *Repository) GetNewsLimited(ctx context.Context, filter map[string]interface{}, limit int64) ([]entity.New, error) {
	r.logger.Debug("fetching limited news", zap.Any("filter", filter), zap.Int64("limit", limit))

	findOptions := options.Find()
	findOptions.SetLimit(limit)

	res, err := r.newsColl.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("failed get limited news", zap.Error(err))
		return nil, fmt.Errorf("get limited news: %w", ErrNotFound)
	}
	defer res.Close(ctx)

	var news []entity.New
	for res.Next(ctx) {
		var n entity.New
		if err := res.Decode(&n); err != nil {
			r.logger.Warn("failed decode new", zap.Error(err))
			return nil, fmt.Errorf("decode news: %w", ErrDecodeFailed)
		}
		news = append(news, n)
	}
	if err = res.Err(); err != nil {
		r.logger.Error("failed parse news", zap.Error(err))
		return nil, fmt.Errorf("parse news: %w", ErrDecodeFailed)
	}

	if len(news) == 0 {
		return nil, ErrNoDocuments
	}

	r.logger.Info("news fetched (limited)", zap.Int("length", len(news)))
	return news, nil
}

func (r *Repository) DeleteNew(ctx context.Context, filter map[string]interface{}) error {
	r.logger.Debug("deleting new", zap.Any("filter", filter))

	res, err := r.newsColl.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("failed delete new", zap.Error(err))
		return fmt.Errorf("delete new: %w", ErrDeleteFailed)
	}

	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	r.logger.Info("new deleted", zap.Any("filter", filter))

	return nil
}
