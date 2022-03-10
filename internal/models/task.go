package models

import (
	"bytes"
	"encoding/json"
	"net/http"
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

// check task validity
func (t *Task) Check() error {
	if err := t.Limits.Check(); err != nil {
		return err
	}

	// set token and default processing result
	t.Token = uuid.New()
	t.Result = Result{
		Token:     t.Token,
		Timestamp: time.Now(),
		Status:    Processing,
	}

	return nil
}

// save result to TaskBucket/{token}
func (t *Task) Save(storage *storage.Client) error {
	config := config.GetConfig()
	object := storage.Bucket(config.TaskBucket).Object(t.Token.String())

	return StorageSave(object, t.Result)
}

// sends result callback to CallbackURL
func (t *Task) SendCallback() {
	if t.CallbackURL == "" {
		return
	}

	body, err := json.Marshal(t.Result)
	if err != nil {
		return
	}

	request, err := http.NewRequest(http.MethodPut, t.CallbackURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	client := &http.Client{
		Timeout: 16 * time.Second,
	}

	// try PUT CallbackURL for at most three times
	for i := 0; i < 3; i++ {
		resp, err := client.Do(request)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
	}
}

// update a task = update timestamp to time.now(), then save the task and send callback
func (t *Task) Update(storage *storage.Client) error {
	t.Result.Timestamp = time.Now()
	err := t.Save(storage)
	t.SendCallback()

	return err
}
