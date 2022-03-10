package runner

import (
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/jason-plainlog/code-exec/internal/models"
)

// each runner handles the whole prepare / compile / execution of all the tasks of a submission
type Runner struct {
	Storage    *storage.Client
	Submission *models.Submission
	TempDir    string
}

// prepare, compile, execute the submission and update the result
func (r *Runner) Handle(s *models.Submission) {
	var err error

	// save submission to struct member
	r.Submission = s
	// create temporary directory to store compiled result
	r.TempDir, err = os.MkdirTemp("", "*")
	if err != nil {
		// update all result status to be Internal Error
		for i := range s.Tasks {
			s.Tasks[i].Result.Status = models.InternalError
			go func(task *models.Task) {
				task.Update(r.Storage)
			}(&s.Tasks[i])
		}
		return
	}
	defer os.RemoveAll(r.TempDir)

	// prepare and compile
	compileResult := r.Compile()
	if compileResult.Status == models.CompileError {
		// update all result status to be Compile Error
		for i := range s.Tasks {
			s.Tasks[i].Result = *compileResult
			s.Tasks[i].Result.Token = s.Tasks[i].Token
			go func(task *models.Task) {
				task.Update(r.Storage)
			}(&s.Tasks[i])
		}
		return
	}

	// run all the tasks parallelly
	for i := range s.Tasks {
		go func(task *models.Task) {
			r.RunTask(task)
			task.Update(r.Storage)
		}(&s.Tasks[i])
	}
}

// compile according to the language configuratino
func (r *Runner) Compile() *models.Result {
	return &models.Result{
		Status:    models.CompileError,
		Timestamp: time.Now(),
	}
}

// run the given task according to the language config
func (r *Runner) RunTask(task *models.Task) {
	task.Result.Status = models.Accepted
	task.Result.Timestamp = time.Now()
}
