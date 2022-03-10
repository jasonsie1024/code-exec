package models

import (
	"fmt"
	"time"

	"cloud.google.com/go/storage"
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

// check submission validity
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

// the submission format stored in storage
type StoredSubmission struct {
	Token     uuid.UUID `json:"token"`
	Timestamp time.Time `json:"timestamp"`

	Langauge   string `json:"langauge"`
	SourceCode []byte `json:"source_code"`

	Tasks []uuid.UUID `json:"tasks"`
}

// save the submission to storage SubmissionBucket/{token}
func (s *Submission) Save(storage *storage.Client) error {
	languages := config.GetLanguages()
	config := config.GetConfig()

	object := storage.Bucket(config.SubmissionBucket).Object(s.Token.String())

	// save tasks first
	errChan := make(chan error)
	for i := range s.Tasks {
		go func(task *Task) {
			errChan <- task.Save(storage)
		}(&s.Tasks[i])
	}
	for range s.Tasks {
		if err := <-errChan; err != nil {
			return err
		}
	}

	// save submission
	store := StoredSubmission{
		Token:      s.Token,
		Timestamp:  s.Timestamp,
		Langauge:   languages[s.LanguageId].Name,
		SourceCode: s.SourceCode,
		Tasks:      []uuid.UUID{},
	}
	for _, task := range s.Tasks {
		store.Tasks = append(store.Tasks, task.Token)
	}

	return StorageSave(object, store)
}
