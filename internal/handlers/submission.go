package handlers

import (
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/jason-plainlog/code-exec/internal/models"
	"github.com/labstack/echo/v4"
)

type SubmissionHandler struct {
	Storage *storage.Client
}

func (h *SubmissionHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/submission", h.Create)
	e.GET("/submission/:token", h.Get)
	e.GET("/task/:token", h.Get)
}

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

func (h *SubmissionHandler) Get(c echo.Context) error {
	token := c.Param("token")

	return c.JSON(http.StatusOK, token)
}

func (h *SubmissionHandler) GetTask(c echo.Context) error {
	token := c.Param("token")

	return c.JSON(http.StatusOK, token)
}
