package handlers

import (
	"context"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/jason-plainlog/code-exec/internal/models"
	"github.com/jason-plainlog/code-exec/internal/runner"
	"github.com/labstack/echo/v4"
)

type SubmissionHandler struct {
	Storage *storage.Client
}

// register handler routes
func (h *SubmissionHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/submission", h.Create)
	e.GET("/submission/:token", h.Get)
	e.GET("/task/:token", h.GetTask)
}

// create submission
func (h *SubmissionHandler) Create(c echo.Context) error {
	// getting submission from request and check validity
	submission := new(models.Submission)
	if err := c.Bind(submission); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err := submission.Check(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// save submission / task to storage
	err := submission.Save(h.Storage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create submission and tasks")
	}

	// call runner to handle submission
	runner := runner.Runner{
		Storage: h.Storage,
	}
	go runner.Handle(submission)

	// generating response
	tasks := []map[string]interface{}{}
	for _, task := range submission.Tasks {
		tasks = append(tasks, map[string]interface{}{
			"token": task.Token,
		})
	}
	response := map[string]interface{}{
		"token": submission.Token,
		"tasks": tasks,
	}

	return c.JSON(http.StatusCreated, response)
}

// get submission
func (h *SubmissionHandler) Get(c echo.Context) error {
	token := c.Param("token")
	config := config.GetConfig()

	// get submission first
	submissionBucket := h.Storage.Bucket(config.SubmissionBucket)
	submission := models.StoredSubmission{}
	err := models.StorageRead(submissionBucket.Object(token), &submission)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "submission not found")
	}

	// get task results
	taskBucket := h.Storage.Bucket(config.TaskBucket)
	resultsChan := make([]chan models.Result, len(submission.Tasks))
	for i, token := range submission.Tasks {
		resultsChan[i] = make(chan models.Result, 1)
		go func(token string, resultChan chan models.Result) {
			result := models.Result{}
			models.StorageRead(taskBucket.Object(token), &result)
			resultChan <- result
		}(token.String(), resultsChan[i])
	}

	// form final response
	tasks := []map[string]interface{}{}
	for i := range submission.Tasks {
		result := <-resultsChan[i]
		tasks = append(tasks, map[string]interface{}{
			"token":  result.Token,
			"time":   result.Time,
			"memory": result.Memory,
			"status": result.Status,
		})
	}
	response := map[string]interface{}{
		"token":       submission.Token,
		"timestamp":   submission.Timestamp,
		"source_code": submission.SourceCode,
		"language":    submission.Langauge,
		"tasks":       tasks,
	}

	return c.JSON(http.StatusOK, response)
}

// get a specific task
func (h *SubmissionHandler) GetTask(c echo.Context) error {
	token := c.Param("token")
	config := config.GetConfig()
	taskBucket := h.Storage.Bucket(config.TaskBucket)

	reader, err := taskBucket.Object(token).NewReader(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "task not found")
	}

	return c.Stream(http.StatusOK, "application/json", reader)
}
