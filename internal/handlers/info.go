package handlers

import (
	"net/http"

	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/labstack/echo/v4"
)

type InfoHandler struct {
}

func (h *InfoHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/config", h.GetConfig)
}

func (h *InfoHandler) GetConfig(c echo.Context) error {
	config := config.GetConfig()
	return c.JSON(http.StatusOK, config)
}
