package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/jason-plainlog/code-exec/internal/config"
)

type (
	Submission struct {
		SourceCode []byte `json:"source_code"`
		Language   string `json:"language"`

		CompileOptions  *string `json:"compile_options"`
		AdditionalFiles []byte  `json:"additional_files"`

		Tasks []Task `json:"tasks"`
	}

	Task struct {
		Stdin  []byte `json:"stdin"`
		Limits Limits `json:"limits"`

		Token uuid.UUID `json:"token"`

		CommandLineArguments *string `json:"command_line_arguments"`
		ExpectedOutput       *string `json:"expected_output"`
		CallbackURL          *string `json:"callback_url"`
	}

	Limits struct {
		Time     float32 `json:"time"`
		Memory   int     `json:"memory"`
		Filesize int     `json:"filesize"`
		Process  int     `json:"process"`
		Network  bool    `json:"network"`
	}

	Result struct {
		Token uuid.UUID `json:"token"`

		Stdout        []byte `json:"stdout"`
		Stderr        []byte `json:"stderr"`
		CompileOutput []byte `json:"compile_output,omitempty"`
		ExitCode      int    `json:"exit_code"`

		Time   float32 `json:"time"`
		Memory int     `json:"memory"`

		Message   string    `json:"message"`
		Status    string    `json:"status"`
		Timestamp time.Time `json:"timestamp"`

		CallbackURL *string `json:"-"`
	}
)

func (s *Submission) Check() error {
	config := config.GetConfig()

	if len(s.SourceCode) == 0 {
		return fmt.Errorf("source_code can't be blank")
	}

	if len(s.Tasks) == 0 || len(s.Tasks) > config.MaxTask {
		return fmt.Errorf("the length of tasks should be at least 1 and at most %d", config.MaxTask)
	}

	for i := 0; i < len(s.Tasks); i++ {
		s.Tasks[i].Check()
	}

	return nil
}

func (t *Task) Check() error {
	if t.Token == uuid.Nil {
		t.Token = uuid.New()
	}

	return t.Limits.Check()
}

func (l *Limits) Check() error {
	config := config.GetConfig()

	if l.Filesize == 0 {
		l.Filesize = config.MaxFilesize
	} else if l.Filesize > config.MaxFilesize {
		return fmt.Errorf("limits.filesize should not exceed %d", config.MaxFilesize)
	}

	if l.Process == 0 {
		l.Process = config.MaxProcess
	} else if l.Process > config.MaxProcess {
		return fmt.Errorf("limits.process should not exceed %d", config.MaxProcess)
	}

	if l.Time == 0 {
		l.Time = config.MaxTime
	} else if l.Time < 0 || l.Time > config.MaxTime {
		return fmt.Errorf("limits.time should be greater than 0.0 and less than %f", config.MaxTime)
	}

	if l.Memory == 0 {
		l.Memory = config.MaxMemory
	} else if l.Memory < 0 || l.Memory > config.MaxMemory {
		return fmt.Errorf("limits.memory should be greater than 0 and less than %d", config.MaxMemory)
	}

	return nil
}

func (r *Result) SendCallback() {
	if r.CallbackURL == nil {
		return
	}

	body, _ := json.Marshal(r)
	go func() {
		client := http.Client{
			Timeout: time.Second * 30,
		}

		// retry at most 3 times
		for i := 0; i < 3; i++ {
			req, err := http.NewRequest(http.MethodPut, *r.CallbackURL, bytes.NewBuffer(body))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil || resp.StatusCode == 200 {
				return
			}
		}
	}()
}

func (r *Result) Save(object *storage.ObjectHandle) {
	writer := object.NewWriter(context.Background())
	defer writer.Close()

	body, _ := json.Marshal(r)
	io.Copy(writer, bytes.NewBuffer(body))
}

var MaximumLimits Limits = Limits{
	Time:     config.GetConfig().MaxTime,
	Memory:   config.GetConfig().MaxMemory,
	Filesize: config.GetConfig().MaxFilesize,
	Process:  config.GetConfig().MaxProcess,
	Network:  true,
}
