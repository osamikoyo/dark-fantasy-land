package repository

import (
	"context"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.uber.org/zap"
)

func (r *Repository) CreateNew(ctx context.Context, New *entity.New) error {
	res, err := r.newsColl.InsertOne(ctx, New)
	if err != nil {
		r.logger.Error("failed create new",
			zap.Error(err))

		return err
	}

	r.logger.Info("new created",
		zap.String("title", New.Title),
		zap.String("inserted_id", res.InsertedID.(string)))

	return nil
}

func (r *Repository) UpdateNew(ctx context.Context, filter map[string]string, update map[string]string) error {
	r.logger.Debug("updating new",
		zap.Any("filter", filter),
		zap.Any("update", update))

	res, err := r.newsColl.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("failed update one", zap.Error(err))
		return err
	}

	r.logger.Info("update one new",
		zap.Int64("matched_count", res.MatchedCount))

	return nil
}

func (r *Repository) GetNews(ctx context.Context, filter map[string]string) ([]entity.New, error) {
	r.logger.Debug("fetching news", zap.Any("filter", filter))

	res, err := r.newsColl.Find(ctx, filter)
	if err != nil {
		r.logger.Error("failed get news", zap.Error(err))

		return nil, err
	}

	defer res.Close(ctx)

	var news []entity.New

	for res.Next(ctx) {
		var new entity.New

		if err := res.Decode(&new); err != nil {
			r.logger.Warn("failed decode new", zap.Error(err))
			continue
		}

		news = append(news, new)
	}

	if err = res.Err(); err != nil {
		r.logger.Error("failed parse news", zap.Error(err))

		return nil, err
	}

	r.logger.Info("news was got",
		zap.Int("length", len(news)))

	return news, nil
}
