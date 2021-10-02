package process

import (
	"log"
	"os/exec"
	"time"
)

type Process struct {
	Command string
	Args    []string
}

type ProcessDaemon struct {
	processes []Process
}

func (pd *ProcessDaemon) Add(p Process) {
	go func() {
		for {
			cmd := exec.Command(p.Command, p.Args...)
			err := cmd.Run()
			log.Printf("Run finish with: %s", err)
			time.Sleep(5 * time.Second)
		}
	}()
}
