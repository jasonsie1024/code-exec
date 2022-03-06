package runners

import (
	"bytes"

	"github.com/jason-plainlog/code-exec/internal/isolate"
	"github.com/jason-plainlog/code-exec/internal/models"
)

type python3Runner struct{}

// handles the whole processing of compiling and executing the submission
func (r *python3Runner) Handle(s *models.Submission) chan *models.Result {
	resultChan := make(chan *models.Result, len(s.Tasks))

	// run each task concurrently
	for i := range s.Tasks {
		go func(i int) {
			result := r.Execute(s, &s.Tasks[i])
			resultChan <- result
		}(i)
	}

	return resultChan
}

func (r *python3Runner) Compile(s *models.Submission) *models.Result {
	// don't need to compile, yah!

	return nil
}

func (r *python3Runner) Execute(s *models.Submission, t *models.Task) *models.Result {
	sandbox := isolate.GetSandbox()
	defer sandbox.CleanUp()

	sandbox.Prepare(s)

	result := sandbox.Run([]string{
		"/usr/bin/python3", "source",
	}, t.Limits, t.Stdin)
	result.CallbackURL = t.CallbackURL
	result.Token = t.Token
	if result.Status == "Accepted" &&
		t.ExpectedOutput != nil &&
		!bytes.Equal(result.Stdout, t.ExpectedOutput) {
		result.Status = "Wrong Answer"
	}

	return result
}
