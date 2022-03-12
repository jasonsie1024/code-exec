package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/jason-plainlog/code-exec/internal/isolate"
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
	if compileResult.Status != models.Accepted {
		// update all result status to be compile result
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
	taskWg := new(sync.WaitGroup)
	for i := range s.Tasks {
		taskWg.Add(1)
		go func(task *models.Task) {
			r.RunTask(task)
			task.Update(r.Storage)
			taskWg.Done()
		}(&s.Tasks[i])
	}
	taskWg.Wait()
}

// compile according to the language configuratino
func (r *Runner) Compile() *models.Result {
	// getting compile stage sandbox
	box, err := isolate.GetSandbox()
	if err != nil {
		return &models.Result{
			Status: models.InternalError,
		}
	}
	defer func() {
		// copy directory to r.tempdir before cleaning up
		exec.Command("cp", "-r", box.Path+"/box", r.TempDir).Run()
		box.CleanUp()
	}()

	// prepare source_code and unzip additional_files
	err = box.Prepare(r.Submission)
	if err != nil {
		return &models.Result{
			Status: models.InternalError,
		}
	}

	lang := config.GetLanguages()[r.Submission.LanguageId]

	result := &models.Result{
		Status: models.Accepted,
	}

	// compilation
	if lang.CompileCommand != "" {
		result = box.Run([]string{
			"/bin/sh", "-c", fmt.Sprintf(lang.CompileCommand, r.Submission.CompileOptions),
		}, models.MaximumLimit, nil)
	}

	if result.Status != models.Accepted {
		result.Status = models.CompileError
	}

	return result
}

// run the given task according to the language config
func (r *Runner) RunTask(task *models.Task) {
	// get task sandbox
	box, err := isolate.GetSandbox()
	if err != nil {
		task.Result.Status = models.InternalError
		task.Update(r.Storage)
		return
	}
	defer box.CleanUp()

	// copy compilation environment
	err = exec.Command("cp", "-r", r.TempDir+"/box", box.Path).Run()
	if err != nil {
		task.Result.Status = models.InternalError
		task.Update(r.Storage)
		return
	}

	lang := config.GetLanguages()[r.Submission.LanguageId]
	// execute task by language config
	result := box.Run([]string{
		"/bin/bash", "-c", fmt.Sprintf(lang.RunCommand, task.CommandLineArguments),
	}, task.Limits, task.Stdin)

	task.Result = *result
	task.Result.Token = task.Token

	if task.Result.Status == models.Accepted && task.ExpectedOutput != nil {
		if !bytes.Equal(task.ExpectedOutput, task.Result.Stdout) {
			task.Result.Status = models.WrongAnswer
		}
	}
}
