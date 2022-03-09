package handlers

import (
	"net/http"
	"strconv"

	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/labstack/echo/v4"
)

type InfoHandler struct {
}

func (h *InfoHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/config", h.GetConfig)
	e.GET("/languages", h.GetLanguages)
	e.GET("/language/:id", h.GetLanguage)
}

func (h *InfoHandler) GetConfig(c echo.Context) error {
	config := config.GetConfig()
	return c.JSON(http.StatusOK, config)
}

func (h *InfoHandler) GetLanguages(c echo.Context) error {
	languages := config.GetLanguages()

	response := []map[string]interface{}{}
	for _, lang := range languages {
		response = append(response, map[string]interface{}{
			"id":   lang.Id,
			"name": lang.Name,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *InfoHandler) GetLanguage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id should be an integer")
	}

	languages := config.GetLanguages()
	lang, found := languages[id]
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "language not exist")
	}

	return echo.NewHTTPError(http.StatusOK, lang)
}
