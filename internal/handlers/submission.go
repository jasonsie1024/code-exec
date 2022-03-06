package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	cfg "github.com/jason-plainlog/code-exec/internal/config"
	"github.com/jason-plainlog/code-exec/internal/models"
	"github.com/jason-plainlog/code-exec/internal/runners"
	"github.com/labstack/echo/v4"

	"cloud.google.com/go/storage"
)

type (
	SubmissionHandler struct {
		StorageClient *storage.Client
	}
)

var config = cfg.GetConfig()

func (h *SubmissionHandler) GetResult(c echo.Context) error {
	token := c.Param("token")
	resultDoc := h.StorageClient.Bucket(config.Bucket).Object(token)

	reader, err := resultDoc.NewReader(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	defer reader.Close()

	return c.Stream(http.StatusOK, "application/json", reader)
}

func (h *SubmissionHandler) Create(c echo.Context) error {
	submission := new(models.Submission)

	// check if the request is a valid submission
	if err := c.Bind(submission); err != nil {
		return err
	} else if err := submission.Check(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	runner, ok := runners.Runners[submission.Language]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid language")
	}

	resultChan := runner.Handle(submission)
	go func() {
		bucket := h.StorageClient.Bucket(config.Bucket)
		for range submission.Tasks {
			result := <-resultChan
			result.SendCallback()
			result.Save(bucket.Object(result.Token.String()))
		}
	}()

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
