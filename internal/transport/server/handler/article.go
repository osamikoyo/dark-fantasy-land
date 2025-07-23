package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/dark-fantasy-land/internal/entity"
)

func (h *Handler) CreateArticle(c echo.Context) error {
	var article entity.Article

	if err := c.Bind(&article); err != nil {
		return c.String(http.StatusBadRequest, ErrInvalidInput)
	}

	article.Timestamp = time.Now()

	if err := h.service.CreateArticle(&article); err != nil {
		return c.String(http.StatusInternalServerError, ErrInternal)
	}

	return c.String(http.StatusCreated, "article created")
}

func (h *Handler) GetArticle(c echo.Context) error {
	author := c.Param("author")
	title := c.Param("title")

	article, err := h.service.GetOneArticle(author, title)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrInternal)
	}

	return c.JSON(http.StatusOK, article)
}

func (h *Handler) GetArticles(c echo.Context) error {
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

	articles, err := h.service.GetMoreArticles(filter)
	if err != nil {
		return c.String(http.StatusInternalServerError, ErrInternal)
	}

	return c.JSON(http.StatusOK, articles)
}
