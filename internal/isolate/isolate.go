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

type (
	Sandbox struct {
		Id   string
		Path string

		inUse bool
	}
)

var availableSandboxes chan *Sandbox

func GetSandbox() *Sandbox {
	if availableSandboxes == nil {
		config := config.GetConfig()
		availableSandboxes = make(chan *Sandbox, config.MaxSandbox)

		go func() {
			for id := 0; id < config.MaxSandbox; id++ {
				availableSandboxes <- &Sandbox{
					Id: fmt.Sprint(id),
				}
			}
		}()
	}

	sandbox := <-availableSandboxes
	for sandbox.init() != nil {
		availableSandboxes <- sandbox
		sandbox = <-availableSandboxes
	}

	return sandbox
}

func (box *Sandbox) init() error {
	err := exec.Command("isolate", "--box-id", box.Id, "--cg", "--cleanup").Run()
	if err != nil {
		return err
	}

	output, err := exec.Command("isolate", "--box-id", box.Id, "--cg", "--init").Output()
	if err != nil {
		return err
	}

	box.Path = strings.TrimSuffix(string(output), "\n")
	box.inUse = true
	return nil
}

func (box *Sandbox) CleanUp() error {
	if !box.inUse {
		return fmt.Errorf("sandbox is already cleaned up")
	}

	if err := exec.Command("isolate", "--box-id", box.Id, "--cg", "--cleanup").Run(); err != nil {
		return err
	}

	box.inUse = false
	availableSandboxes <- box
	return nil
}

func (box *Sandbox) Run(command []string, limits models.Limits, stdin []byte) *models.Result {
	// build up command: `isolate {boxid, limits} --run {command, arguments}`
	args := []string{
		"--box-id", box.Id, "--cg",
		"-M", box.Path + "/meta",
		"-t", fmt.Sprint(limits.Time), "-w", fmt.Sprint(limits.Time * 2),
		"-m", fmt.Sprint(limits.Memory), "-f", fmt.Sprint(limits.Filesize),
		fmt.Sprintf("--processe=%d", limits.Process),
		"-o", "stdout.txt", "-r", "stderr.txt",
	}
	if limits.Network {
		args = append(args, "--share-net")
	}
	args = append(args, "--run", "--")
	args = append(args, command...)

	// execute and write stdin if exist
	cmd := exec.Command("isolate", args...)
	if stdin != nil {
		Stdin, _ := cmd.StdinPipe()
		Stdin.Write(stdin)
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(output), err)
	}
	meta := ParseMetafile(box.Path + "/meta")

	result := &models.Result{
		Status:    meta["status"],
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

func (box *Sandbox) Prepare(s *models.Submission) {
	os.WriteFile(box.Path+"/box/source", s.SourceCode, 0644)

	if s.AdditionalFiles != nil {
		// dump zip file
		os.WriteFile(box.Path+"/box/files.zip", s.AdditionalFiles, 0644)

		// safely unzip
		box.Run([]string{
			"/usr/bin/unzip", "-n", "files.zip",
		}, models.MaximumLimits, nil)

		// delete zip file
		os.Remove(box.Path + "/box/files.zip")
	}
}
