package handlers

import (
	"net/http"

	"github.com/jason-plainlog/code-exec/internal/models"
	"github.com/labstack/echo/v4"
)

type (
	SubmissionHandler struct {
	}
)

func (h *SubmissionHandler) Get(c echo.Context) error {

	return nil
}

func (h *SubmissionHandler) Create(c echo.Context) error {
	submission := new(models.Submission)
	if err := c.Bind(submission); err != nil {
		return err
	}

	if err := submission.Check(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, submission)
}
