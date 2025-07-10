package repository

import (
	"context"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (r *Repository) CreateMem(ctx context.Context, mem *entity.Mem) error {
	res, err := r.cfuColl.InsertOne(ctx, mem)
	if err != nil {
		r.logger.Error("failed to create mem", zap.Error(err))
		return err
	}

	r.logger.Info("mem created",
		zap.String("image_name", mem.ImageName),
		zap.String("inserted_id", res.InsertedID.(string)))
	return nil
}

func (r *Repository) UpdateMem(ctx context.Context, filter, update map[string]interface{}) error {
	r.logger.Debug("updating mem", zap.Any("filter", filter), zap.Any("update", update))
	res, err := r.cfuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("failed to update mem", zap.Error(err))
		return err
	}

	r.logger.Info("mem updated", zap.Int64("matched_count", res.MatchedCount))
	return nil
}

func (r *Repository) GetMems(ctx context.Context, filter map[string]interface{}) ([]entity.Mem, error) {
	r.logger.Debug("fetching mems", zap.Any("filter", filter))
	res, err := r.cfuColl.Find(ctx, filter)
	if err != nil {
		r.logger.Error("failed to get mems", zap.Error(err))
		return nil, err
	}
	defer res.Close(ctx)

	var mems []entity.Mem
	for res.Next(ctx) {
		var mem entity.Mem
		if err := res.Decode(&mem); err != nil {
			r.logger.Warn("failed to decode mem", zap.Error(err))
			continue
		}
		mems = append(mems, mem)
	}
	if err = res.Err(); err != nil {
		r.logger.Error("failed to parse mems", zap.Error(err))
		return nil, err
	}

	r.logger.Info("mems fetched", zap.Int("length", len(mems)))
	return mems, nil
}

func (r *Repository) GetMemsLimited(ctx context.Context, filter map[string]interface{}, limit int64) ([]entity.Mem, error) {
	r.logger.Debug("fetching limited mems", zap.Any("filter", filter), zap.Int64("limit", limit))

	findOptions := options.Find()
	findOptions.SetLimit(limit)

	res, err := r.cfuColl.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("failed to get limited mems", zap.Error(err))
		return nil, err
	}
	defer res.Close(ctx)

	var mems []entity.Mem
	for res.Next(ctx) {
		var mem entity.Mem
		if err := res.Decode(&mem); err != nil {
			r.logger.Warn("failed to decode mem", zap.Error(err))
			continue
		}
		mems = append(mems, mem)
	}
	if err = res.Err(); err != nil {
		r.logger.Error("failed to parse mems", zap.Error(err))
		return nil, err
	}

	r.logger.Info("mems fetched (limited)", zap.Int("length", len(mems)))
	return mems, nil
}

func (r *Repository) DeleteMem(ctx context.Context, filter map[string]interface{}) error {
	r.logger.Debug("deleting mem", zap.Any("filter", filter))
	_, err := r.cfuColl.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("failed to delete mem", zap.Error(err))
		return err
	}

	r.logger.Info("mem deleted", zap.Any("filter", filter))
	return nil
}
