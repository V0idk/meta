package error

import (
	"fmt"
	pb "meta/msg"
)

type MSGTYPE_NOT_FOUND struct{}

func (e MSGTYPE_NOT_FOUND) Error() string {
	return "MSGTYPE_NOT_FOUND"
}

type PROCESS_NOT_FOUND struct{}

func (e PROCESS_NOT_FOUND) Error() string {
	return "PROCESS_NOT_FOUND"
}

type ENTRY_NOT_FOUND struct{}

func (e ENTRY_NOT_FOUND) Error() string {
	return "ENTRY_NOT_FOUND"
}

type PARAM_WRONG struct {
	Param string
}

func (e PARAM_WRONG) Error() string {
	return fmt.Sprintf("PARAM_WRONG: %s", e.Param)
}

func GetErrorMsg(e error) (*pb.Msg, error) {
	return &pb.Msg{
		Type:    pb.ERR.Id,
		Content: []byte(e.Error()),
	}, e
}
