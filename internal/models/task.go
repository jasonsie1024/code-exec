package models

import (
	"time"

	"github.com/google/uuid"
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
