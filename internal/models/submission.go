package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jason-plainlog/code-exec/internal/config"
)

type (
	Submission struct {
		Token     uuid.UUID `json:"token"`
		Timestamp time.Time `json:"timestamp"`

		// required
		LanguageId int    `json:"language_id"`
		SourceCode []byte `json:"source_code"`
		Tasks      []Task `json:"tasks"`

		// optional
		CompileOptions  string `json:"compile_options"`
		AdditionalFiles []byte `json:"additional_files"`
	}
)

func (s *Submission) Check() error {
	if s.SourceCode == nil {
		return fmt.Errorf("source_code is required")
	}

	languages := config.GetLanguages()
	if _, exist := languages[s.LanguageId]; !exist {
		return fmt.Errorf("language doesn't exist")
	}

	config := config.GetConfig()
	if len(s.Tasks) == 0 || len(s.Tasks) > config.MaxTask {
		return fmt.Errorf("the length of tasks should be at least 1 and at most %d", config.MaxTask)
	}

	for i := range s.Tasks {
		if err := s.Tasks[i].Check(); err != nil {
			return err
		}
	}

	s.Token = uuid.New()
	s.Timestamp = time.Now()

	return nil
}
