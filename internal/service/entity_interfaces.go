package service

import (
	"context"

	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
)

type (
	Repository interface {
		ArticleRepository
		NewRepository
		WallpaperRepository
		MemRepository
	}

	Casher interface {
		ArticleCasher
		NewCasher
		MemCasher
		WallpaperCasher
	}

	Sender interface {
		SendToCensor(string, interface{}) error
	}

	ArticleCasher interface {
		AddArticleToCash(context.Context, *entity.Article) error
		UpdateArticleInCash(context.Context, string, string, string, interface{}) error
		GetArticleFromCash(context.Context, string, string) (*entity.Article, error)
		DeleteArticleFromCash(context.Context, string, string) error
	}

	MemCasher interface {
		AddMemToCash(context.Context, *entity.Mem) error
		UpdateMemInCash(context.Context, string, string, string, interface{}) error
		GetMemFromCash(context.Context, string, string) (*entity.Mem, error)
		DeleteMemInCash(context.Context, string, string) error
	}

	NewCasher interface {
		AddNewToCash(context.Context, *entity.New) error
		UpdateNewInCash(context.Context, string, string, string, interface{}) error
		GetNewFromCash(context.Context, string, string) (*entity.New, error)
		DeleteNewFromCash(context.Context, string, string) error
	}

	WallpaperCasher interface {
		AddWallpaperToCash(context.Context, *entity.Wallpaper) error
		UpdateWallpaperInCash(context.Context, string, string, string, interface{}) error
		GetWallpaperFromCash(context.Context, string, string) (*entity.Wallpaper, error)
		DeleteWallpaperFromCash(context.Context, string, string) error
	}

	ArticleRepository interface {
		CreateArticle(context.Context, *entity.Article) error
		UpdateArticle(context.Context, map[string]interface{}, map[string]interface{}) error
		DeleteArticle(context.Context, map[string]interface{}) error
		GetArticle(context.Context, map[string]interface{}) (*entity.Article, error)
		GetArticlesLimited(context.Context, map[string]interface{}, int64) ([]entity.Article, error)
	}

	MemRepository interface {
		CreateMem(context.Context, *entity.Mem) error
		UpdateMem(context.Context, map[string]interface{}, map[string]interface{}) error
		DeleteMem(context.Context, map[string]interface{}) error
		GetMem(context.Context, map[string]interface{}) (*entity.Mem, error)
		GetMemsLimited(context.Context, map[string]interface{}, int64) ([]entity.Mem, error)
	}

	NewRepository interface {
		CreateNew(context.Context, *entity.New) error
		UpdateNew(context.Context, map[string]interface{}, map[string]interface{}) error
		DeleteNew(context.Context, map[string]interface{}) error
		GetNew(context.Context, map[string]interface{}) (*entity.New, error)
		GetNewsLimited(context.Context, map[string]interface{}, int64) ([]entity.New, error)
	}

	WallpaperRepository interface {
		CreateWallpaper(context.Context, *entity.Wallpaper) error
		UpdateWallpaper(context.Context, map[string]interface{}, map[string]interface{}) error
		DeleteWallpaper(context.Context, map[string]interface{}) error
		GetWallpaper(context.Context, map[string]interface{}) (*entity.Wallpaper, error)
		GetWallpapersLimited(context.Context, map[string]interface{}, int64) ([]entity.Wallpaper, error)
	}
)
