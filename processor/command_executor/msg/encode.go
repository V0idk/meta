package msg

import "bytes"

type CommandContent struct {
	Command string `json:"command"`
	Args    string `json:"args"`
}

type CommandResult struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
	Error  error
}
