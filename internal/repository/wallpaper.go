package repository

import (
	"context"
	"fmt"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) CreateWallpaper(ctx context.Context, wallpaper *entity.Wallpaper) error {
	res, err := r.wallpaperColl.InsertOne(ctx, wallpaper)
	if err != nil {
		r.logger.Error("failed to create wallpaper", zap.Error(err))
		return fmt.Errorf("create wallpaper: %w", ErrInsertFailed)
	}

	r.logger.Info("wallpaper created",
		zap.String("image_name", wallpaper.ImageName),
		zap.String("inserted_id", fmt.Sprintf("%v", res.InsertedID)))
	return nil
}

func (r *Repository) UpdateWallpaper(ctx context.Context, filter, update map[string]interface{}) error {
	r.logger.Debug("updating wallpaper", zap.Any("filter", filter), zap.Any("update", update))
	res, err := r.wallpaperColl.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("failed to update wallpaper", zap.Error(err))
		return fmt.Errorf("update wallpaper: %w", ErrUpdateFailed)
	}

	if res.MatchedCount == 0 {
		return ErrNotFound
	}

	r.logger.Info("wallpaper updated", zap.Int64("matched_count", res.MatchedCount))
	return nil
}

func (r *Repository) GetWallpaper(ctx context.Context, filter map[string]interface{}) (*entity.Wallpaper, error) {
	r.logger.Debug("fetching single wallpaper", zap.Any("filter", filter))

	res := r.wallpaperColl.FindOne(ctx, filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			r.logger.Warn("wallpaper not found", zap.Any("filter", filter))
			return nil, ErrNotFound
		}
		r.logger.Error("failed to get wallpaper", zap.Error(res.Err()))
		return nil, fmt.Errorf("get wallpaper: %w", res.Err())
	}

	var wallpaper entity.Wallpaper
	if err := res.Decode(&wallpaper); err != nil {
		r.logger.Warn("failed to decode wallpaper", zap.Error(err))
		return nil, fmt.Errorf("decode wallpaper: %w", ErrDecodeFailed)
	}

	r.logger.Info("wallpaper fetched", zap.Any("wallpaper", wallpaper))
	return &wallpaper, nil
}

func (r *Repository) GetWallpapersLimited(ctx context.Context, filter map[string]interface{}, limit int64) ([]entity.Wallpaper, error) {
	r.logger.Debug("fetching limited wallpapers", zap.Any("filter", filter), zap.Int64("limit", limit))

	findOptions := options.Find()
	findOptions.SetLimit(limit)

	res, err := r.wallpaperColl.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("failed to get limited wallpapers", zap.Error(err))
		return nil, fmt.Errorf("get limited wallpapers: %w", ErrNotFound)
	}
	defer res.Close(ctx)

	var wallpapers []entity.Wallpaper
	for res.Next(ctx) {
		var wallpaper entity.Wallpaper
		if err := res.Decode(&wallpaper); err != nil {
			r.logger.Warn("failed to decode wallpaper", zap.Error(err))
			return nil, fmt.Errorf("decode wallpaper: %w", ErrDecodeFailed)
		}
		wallpapers = append(wallpapers, wallpaper)
	}
	if err = res.Err(); err != nil {
		r.logger.Error("failed to parse wallpapers", zap.Error(err))
		return nil, fmt.Errorf("parse wallpapers: %w", ErrDecodeFailed)
	}

	if len(wallpapers) == 0 {
		return nil, ErrNoDocuments
	}

	r.logger.Info("wallpapers fetched (limited)", zap.Int("length", len(wallpapers)))
	return wallpapers, nil
}

func (r *Repository) DeleteWallpaper(ctx context.Context, filter map[string]interface{}) error {
	r.logger.Debug("deleting wallpaper", zap.Any("filter", filter))
	res, err := r.wallpaperColl.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("failed to delete wallpaper", zap.Error(err))
		return fmt.Errorf("delete wallpaper: %w", ErrDeleteFailed)
	}

	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	r.logger.Info("wallpaper deleted", zap.Any("filter", filter))
	return nil
}
