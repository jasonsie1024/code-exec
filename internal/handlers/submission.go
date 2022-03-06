package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jason-plainlog/code-exec/internal/models"
	"github.com/labstack/echo/v4"
)

type (
	SubmissionHandler struct {
	}
)

func (h *SubmissionHandler) GetResult(c echo.Context) error {

	return nil
}

func (h *SubmissionHandler) Create(c echo.Context) error {
	submission := new(models.Submission)

	// check if the request is a valid submission
	if err := c.Bind(submission); err != nil {
		return err
	} else if err := submission.Check(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tokens := []uuid.UUID{}
	for _, task := range submission.Tasks {
		tokens = append(tokens, task.Token)
	}

	return c.JSON(http.StatusOK, struct {
		Tokens []uuid.UUID `json:"tokens"`
	}{
		Tokens: tokens,
	})
}
