package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
)

func (h *Handler) CreateMem(c echo.Context) error {
	var mem entity.Mem

	if err := c.Bind(mem); err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidInput)
	}

	if err := h.service.CreateMem(&mem); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "mem created")
}

func (h *Handler) GetMemInfo(c echo.Context) error {
	image_name := c.Param("image_name")
	author := c.Param("author")

	mem, err := h.service.GetOneMem(image_name, author)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, mem)
}

func (h *Handler) GetMemImage(c echo.Context) error {
	image_name := c.Param("image_name")

	obj, err := h.storage.DownloadFile(image_name, "mems")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Stream(http.StatusOK, "application/octet-stream", obj)
}

func (h *Handler) GetMems(c echo.Context) error {
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
