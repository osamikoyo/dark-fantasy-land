package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/dark-fantasy-land/internal/config"
	"github.com/osamikoyo/dark-fantasy-land/internal/service"
	"github.com/osamikoyo/dark-fantasy-land/pkg/storage"
)

type Handler struct {
	service *service.Service
	storage *storage.Storage

	cfg *config.Config
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRouters(e *echo.Echo) {
	articles := e.Group("/article")

	articles.POST("/create", h.CreateArticle)
	articles.GET("/get/one", h.GetArticle)
	articles.GET("/get/more", h.GetArticles)

	mems := e.Group("/mem")

	mems.POST("/create", h.CreateMem)
	mems.GET("/get/info", h.GetMemInfo)
	mems.GET("/get/image", h.GetMemImage)
	mems.GET("/get/more", h.GetMems)

	wallpapers := e.Group("/wallpaper")

	wallpapers.POST("/create", h.CreateWallpaper)
	wallpapers.GET("/get/info", h.GetWallpaperInfo)
	wallpapers.GET("/get/image", h.GetWallpaperImage)
	wallpapers.GET("/get/more", h.GetWallpapers)

	news := e.Group("/news")

	news.POST("/create", h.CreateArticle)
	news.GET("/get/one", h.GetNew)
	news.GET("/get/more", h.GetNews)
}
