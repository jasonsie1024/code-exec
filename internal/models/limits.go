package models

import (
	"fmt"

	"github.com/jason-plainlog/code-exec/internal/config"
)

type Limits struct {
	Time     float32 `json:"time"`
	Memory   int     `json:"memory"`
	Filesize int     `json:"filesize"`
	Process  int     `json:"process"`
	Network  bool    `json:"network"`
}

var MaximumLimit = Limits{
	Time:     config.GetConfig().MaxTime,
	Memory:   config.GetConfig().MaxMemory,
	Filesize: config.GetConfig().MaxFilesize,
	Process:  config.GetConfig().MaxProcess,
	Network:  true,
}

// check limits validity
func (l *Limits) Check() error {
	config := config.GetConfig()

	if l.Time == 0 {
		l.Time = config.MaxTime
	} else if l.Time < 0 || l.Time > config.MaxTime {
		return fmt.Errorf("the time limit should be greater than 0 and less than %f", config.MaxTime)
	}

	if l.Memory == 0 {
		l.Memory = config.MaxMemory
	} else if l.Memory < 0 || l.Memory > config.MaxMemory {
		return fmt.Errorf("the memory limit should be greater than 0 and less than %d", config.MaxMemory)
	}

	if l.Filesize == 0 {
		l.Filesize = config.MaxFilesize
	} else if l.Filesize < 0 || l.Filesize > config.MaxFilesize {
		return fmt.Errorf("the filesize limit should be greater than 0 and less than %d", config.MaxFilesize)
	}

	if l.Process == 0 {
		l.Process = 1
	} else if l.Process < 0 || l.Process > config.MaxProcess {
		return fmt.Errorf("the process limit should be greater than 0 and less than %d", config.MaxProcess)
	}

	return nil
}
