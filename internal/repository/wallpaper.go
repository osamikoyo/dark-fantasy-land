package repository

import (
	"context"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.uber.org/zap"
)

func (r *Repository) CreateWallpaper(ctx context.Context, wallpaper *entity.Wallpaper) error {
	res, err := r.wallpaperColl.InsertOne(ctx, wallpaper)
	if err != nil {
		r.logger.Error("failed to create wallpaper", zap.Error(err))
		return err
	}

	r.logger.Info("wallpaper created",
		zap.String("image_name", wallpaper.ImageName),
		zap.String("inserted_id", res.InsertedID.(string)))
	return nil
}

func (r *Repository) UpdateWallpaper(ctx context.Context, filter map[string]string, update map[string]string) error {
	r.logger.Debug("updating wallpaper", zap.Any("filter", filter), zap.Any("update", update))
	res, err := r.wallpaperColl.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("failed to update wallpaper", zap.Error(err))
		return err
	}

	r.logger.Info("wallpaper updated", zap.Int64("matched_count", res.MatchedCount))
	return nil
}

func (r *Repository) GetWallpapers(ctx context.Context, filter map[string]string) ([]entity.Wallpaper, error) {
	r.logger.Debug("fetching wallpapers", zap.Any("filter", filter))
	res, err := r.wallpaperColl.Find(ctx, filter)
	if err != nil {
		r.logger.Error("failed to get wallpapers", zap.Error(err))
		return nil, err
	}
	defer res.Close(ctx)

	var wallpapers []entity.Wallpaper
	for res.Next(ctx) {
		var wallpaper entity.Wallpaper
		if err := res.Decode(&wallpaper); err != nil {
			r.logger.Warn("failed to decode wallpaper", zap.Error(err))
			continue
		}
		wallpapers = append(wallpapers, wallpaper)
	}
	if err = res.Err(); err != nil {
		r.logger.Error("failed to parse wallpapers", zap.Error(err))
		return nil, err
	}

	r.logger.Info("wallpapers fetched", zap.Int("length", len(wallpapers)))
	return wallpapers, nil
}

func (r *Repository) DeleteWallpaper(ctx context.Context, filter map[string]string) error {
	r.logger.Debug("deleting wallpaper", zap.Any("filter", filter))
	_, err := r.wallpaperColl.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("failed to delete wallpaper", zap.Error(err))
		return err
	}

	r.logger.Info("wallpaper deleted", zap.Any("filter", filter))
	return nil
}
