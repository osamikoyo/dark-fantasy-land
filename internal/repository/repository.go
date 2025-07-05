package repository

import (
	"errors"

	"github.com/osamikoyo/dark-fantasy-land/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	articlesColl  *mongo.Collection
	newsColl      *mongo.Collection
	cfuColl       *mongo.Collection
	wallpaperColl *mongo.Collection
	logger        *logger.Logger
}

func NewRepository(db *mongo.Database, logger *logger.Logger) (*Repository, error) {
	articles := db.Collection("articles")
	if articles == nil {
		return nil, errors.New("failed get collection for articles")
	}

	news := db.Collection("news")
	if news == nil {
		return nil, errors.New("failed get collection for news")
	}

	cfu := db.Collection("cfu")
	if cfu == nil {
		return nil, errors.New("failed get collection for cfu")
	}

	wallpaper := db.Collection("wallpaper")
	if wallpaper == nil {
		return nil, errors.New("failed get collection for wallpaper")
	}

	return &Repository{
		articlesColl:  articles,
		newsColl:      news,
		cfuColl:       cfu,
		wallpaperColl: wallpaper,
		logger:        logger,
	}, nil
}
