package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	InfoHandler struct {
	}
)

func (h *InfoHandler) ConfigInfo(c echo.Context) error {

	return c.JSON(http.StatusOK, config)
}
