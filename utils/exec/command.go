package exec

import (
	"bytes"
	"log"
	"os/exec"
	"syscall"
)

func ExecCommand(name string, arg ...string) ([]byte, []byte, int) {
	var err error
	var exitCode int
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			log.Printf("Failed to get exit code, set default")
			exitCode = 1
			if _, err1 := cmd.Stderr.Write([]byte(err.Error())); err1 != nil {
				log.Printf("Failed to Write")
			}
		}
	} else {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			exitCode = ws.ExitStatus()
		} else {
			log.Printf("Failed to get exit code, set default")
			exitCode = 0
		}
	}
	return stdout.Bytes(), stderr.Bytes(), exitCode
}
