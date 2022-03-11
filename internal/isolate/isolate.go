package isolate

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jason-plainlog/code-exec/internal/config"
	"github.com/jason-plainlog/code-exec/internal/models"
)

type Sandbox struct {
	Id          string
	Path        string
	initialized bool
}

var availableSandbox chan *Sandbox

// should be called once when server starts
func Init() {
	if availableSandbox != nil {
		return
	}

	config := config.GetConfig()
	availableSandbox = make(chan *Sandbox, config.MaxSandbox)
	go func() {
		for i := 0; i < config.MaxSandbox; i++ {
			availableSandbox <- &Sandbox{
				Id:          fmt.Sprint(i),
				initialized: false,
			}
		}
	}()
}

func GetSandbox() (*Sandbox, error) {
	sandbox := <-availableSandbox

	// clean up first
	exec.Command("isolate", "--cg", "--box-id", sandbox.Id, "--cleanup").Run()

	// initialize isolate sandbox
	output, err := exec.Command("isolate", "--cg", "--box-id", sandbox.Id, "--init").Output()
	if err != nil {
		availableSandbox <- sandbox
		return nil, err
	}

	sandbox.Path = strings.TrimSuffix(string(output), "\n")
	sandbox.initialized = true

	return sandbox, nil
}

func (box *Sandbox) Prepare(submission *models.Submission) error {
	if submission.AdditionalFiles != nil {
		// write zip file
		err := os.WriteFile(box.Path+"/box/files.zip", submission.AdditionalFiles, 0644)
		if err != nil {
			return err
		}

		// unzip with sandbox
		result := box.Run([]string{"/usr/bin/unzip", "files.zip"}, models.MaximumLimit, nil)
		if result.Status != models.Accepted {
			return fmt.Errorf("failed to unzip additional_files")
		}

		// remove zip file
		os.Remove(box.Path + "/box/files.zip")
	}

	languages := config.GetLanguages()
	lang := languages[submission.LanguageId]
	err := os.WriteFile(box.Path+"/box/"+lang.SourceFile, submission.SourceCode, 0644)

	return err
}

func (box *Sandbox) Run(cmd []string, limits models.Limits, stdin []byte) *models.Result {
	// delete stdout.txt and stderr.txt if exist
	os.Remove(box.Path + "/box/stdout.txt")
	os.Remove(box.Path + "/box/stderr.txt")

	// build up command: `isolate {boxid, limits} --run {command, arguments}`
	args := []string{
		"--box-id", box.Id, "--cg",
		"-M", box.Path + "/meta",
		"-t", fmt.Sprint(limits.Time), "-w", fmt.Sprint(limits.Time * 2),
		"-m", fmt.Sprint(limits.Memory), "-f", fmt.Sprint(limits.Filesize),
		fmt.Sprintf("--processe=%d", limits.Process),
		"-o", "stdout.txt", "-r", "stderr.txt",
		"-E", "PATH", "-d", "/var", "-d", "/etc:noexec",
	}
	if limits.Network {
		args = append(args, "--share-net")
	}
	args = append(args, "--run", "--")
	args = append(args, cmd...)

	// execute and write stdin if exist
	command := exec.Command("isolate", args...)
	if stdin != nil {
		Stdin, _ := command.StdinPipe()
		Stdin.Write(stdin)
	}

	if output, err := command.CombinedOutput(); err != nil {
		fmt.Println(string(output), err)
	}
	meta := ParseMetafile(box.Path + "/meta")

	result := &models.Result{
		Status:    models.Status(meta["status"]),
		Message:   meta["message"],
		Timestamp: time.Now(),
	}

	result.Stdout, _ = os.ReadFile(box.Path + "/box/stdout.txt")
	result.Stderr, _ = os.ReadFile(box.Path + "/box/stderr.txt")
	fmt.Sscanf(meta["time"], "%f", &result.Time)
	fmt.Sscanf(meta["cg-mem"], "%d", &result.Memory)
	fmt.Sscanf(meta["exitcode"], "%d", &result.ExitCode)

	return result
}

func (box *Sandbox) CleanUp() error {
	if !box.initialized {
		return fmt.Errorf("sandbox is already cleaned up")
	}

	err := exec.Command("isolate", "--cg", "--box-id", box.Id, "--cleanup").Run()
	if err != nil {
		return err
	}
	box.initialized = false
	availableSandbox <- box

	return nil
}
