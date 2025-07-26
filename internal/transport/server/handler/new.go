package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
)

func (h *Handler) CreateNew(c echo.Context) error {
	var new entity.New

	if err := c.Bind(&new); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := h.service.CreateNew(&new); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "new created")
}

func (h *Handler) GetNew(c echo.Context) error {
	author := c.Param("author")
	title := c.Param("title")

	new, err := h.service.GetOneNew(author, title)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, new)
}

func (h *Handler) GetNews(c echo.Context) error {
	params := []string{"author", "title", "timestamp", "content"}

	filter := make(map[string]interface{})
	for _, p := range params {
		filter[p] = c.Param(p)
	}

	for key, value := range filter {
		if value == "" {
			delete(filter, key)
		}
	}

	news, err := h.service.GetManyNew(filter)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, news)
}
