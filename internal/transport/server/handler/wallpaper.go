package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
)

func (h *Handler) CreateWallpaper(c echo.Context) error {
	var wallpaper entity.Wallpaper

	if err := c.Bind(&wallpaper); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err = h.storage.UploadFile(file, h.cfg.MinioBuckets.WallpaperFull); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if err = h.storage.UploadAndCommpress(file, h.cfg.MinioBuckets.WallpaperWatch); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if err = h.service.CreateWallpaper(&wallpaper); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "wallpaper created")
}

func (h *Handler) GetWallpapers(c echo.Context) error {
	params := []string{"author", "title", "image_name", "timestamp"}

	filter := make(map[string]interface{})
	for _, p := range params {
		filter[p] = c.Param(p)
	}

	for key, value := range filter {
		if value == "" {
			delete(filter, key)
		}
	}

	wallpapers, err := h.service.GetManyWallpapers(filter)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, wallpapers)
}

func (h *Handler) GetWallpaperInfo(c echo.Context) error {
	title := c.Param("title")
	image_name := c.Param("image_name")

	wallpaper, err := h.service.GetOneWallpaper(image_name, title)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, wallpaper)
}

func (h *Handler) GetWallpaperImage(c echo.Context) error {
	image_name := c.Param("image_name")

	wallpaper, err := h.storage.DownloadFile(image_name, h.cfg.MinioBuckets.WallpaperWatch)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Stream(http.StatusOK, "application/octet-stream", wallpaper)
}

func (h *Handler) DownloadWallpaper(c echo.Context) error {
	image_name := c.Param("image_name")

	wallpaper, err := h.storage.DownloadFile(image_name, h.cfg.MinioBuckets.WallpaperFull)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Stream(http.StatusOK, "application/octet-stream", wallpaper)
}
