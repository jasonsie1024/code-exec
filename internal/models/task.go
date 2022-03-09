package models

import (
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/jason-plainlog/code-exec/internal/config"
)

type (
	Task struct {
		Token uuid.UUID `json:"token"`

		// request
		Stdin                []byte `json:"stdin"`
		Limits               Limits `json:"limits"`
		ExpectedOutput       []byte `json:"expected_output"`
		CommandLineArguments string `json:"command_line_arguments"`
		CallbackURL          string `json:"callback_url"`

		Result Result `json:"result"`
	}
)

func (t *Task) Check() error {
	if err := t.Limits.Check(); err != nil {
		return err
	}

	t.Token = uuid.New()
	t.Result = Result{
		Token:     t.Token,
		Timestamp: time.Now(),
		Status:    Processing,
	}

	return nil
}

func (t *Task) Save(storage *storage.Client) error {
	config := config.GetConfig()
	object := storage.Bucket(config.TaskBucket).Object(t.Token.String())

	return StorageSave(object, t.Result)
}
