package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/dark-fantasy-land/internal/service"
	"github.com/osamikoyo/dark-fantasy-land/pkg/storage"
)

var (
	ErrInvalidInput = "invalid input"
	ErrInternal     = "internal server error"
)

type Handler struct {
	service *service.Service
	storage *storage.Storage
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRouters(e *echo.Echo) {
	articles := e.Group("/articles")

	articles.GET("/create", h.CreateArticle)
	articles.GET("/get", h.GetArticle)
	articles.GET("/getmany", h.GetArticles)
}
