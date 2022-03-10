package models

import (
	"time"

	"github.com/google/uuid"
)

type Result struct {
	Token     uuid.UUID `json:"token"`
	Timestamp time.Time `json:"timestamp"`

	Stdout []byte `json:"stdout"`
	Stderr []byte `json:"stderr"`

	Time   float32 `json:"time"`
	Memory int     `json:"memory"`

	Status   Status `json:"status"`
	Message  string `json:"message"`
	ExitCode int    `json:"exit_code"`
}

type Status string

const (
	Processing        Status = "Processing"
	Accepted          Status = "Accepted"
	WrongAnswer       Status = "Wrong Answer"
	CompileError      Status = "Compile Error"
	RuntimeError      Status = "Runtime Error"
	TimeLimitExceeded Status = "Time Limit Exceeded"
	InternalError     Status = "Internal Error"
)
