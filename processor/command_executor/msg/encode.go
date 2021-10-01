package msg

type CommandContent struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type CommandResult struct {
	Stdout []byte
	Stderr []byte
	Status int
}
