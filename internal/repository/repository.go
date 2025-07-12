package repository

import (
	"errors"
	"fmt"

	"github.com/osamikoyo/dark-fantasy-land/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("already exists")
	ErrDBConnection  = errors.New("database connection error")
	ErrInsertFailed  = errors.New("insert failed")
	ErrUpdateFailed  = errors.New("update failed")
	ErrDeleteFailed  = errors.New("delete failed")
	ErrDecodeFailed  = errors.New("decode failed")
	ErrNoDocuments   = errors.New("no documents in result")
)

type Repository struct {
	articlesColl  *mongo.Collection
	newsColl      *mongo.Collection
	cfuColl       *mongo.Collection
	wallpaperColl *mongo.Collection
	userColl      *mongo.Collection
	logger        *logger.Logger
}

func NewRepository(db *mongo.Database, logger *logger.Logger) (*Repository, error) {
	articles := db.Collection("articles")
	if articles == nil {
		return nil, fmt.Errorf("failed get collection for articles: %w", ErrNotFound)
	}

	news := db.Collection("news")
	if news == nil {
		return nil, fmt.Errorf("failed get collection for news: %w", ErrNotFound)
	}

	cfu := db.Collection("cfu")
	if cfu == nil {
		return nil, fmt.Errorf("failed get collection for cfu: %w", ErrNotFound)
	}

	wallpaper := db.Collection("wallpaper")
	if wallpaper == nil {
		return nil, fmt.Errorf("failed get collection for wallpaper: %w", ErrNotFound)
	}

	return &Repository{
		articlesColl:  articles,
		newsColl:      news,
		cfuColl:       cfu,
		wallpaperColl: wallpaper,
		logger:        logger,
	}, nil
}
