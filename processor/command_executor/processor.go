package command_executor

import (
	"encoding/json"
	"log"
	self_error "meta/error"
	pb "meta/msg"
	. "meta/processor/command_executor/msg"
	"os/exec"
)

type CommandExecutor struct {
}

func (c *CommandExecutor) Command(in *pb.Msg) (*pb.Msg, error) {
	var content CommandContent
	err := json.Unmarshal(in.Content, &content)
	if err != nil {
		log.Printf("Failed to load")
		return nil, err
	}
	result := CommandResult{}
	cmd := exec.Command(content.Command, content.Args)
	cmd.Stdout = &result.Stdout
	cmd.Stderr = &result.Stderr
	result.Error = cmd.Run()

	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Printf("could not marshal: %s", result)
		return nil, err
	}
	return &pb.Msg{
		Type:    pb.OK.Id,
		Content: resultBytes,
	}, nil
}

func (c *CommandExecutor) Dispatch(in *pb.Msg) (*pb.Msg, error) {
	switch in.Type {
	case pb.COMMAND.Id:
		return c.Command(in)
	}
	return nil, &self_error.MSGTYPE_NOT_FOUND{}
}
