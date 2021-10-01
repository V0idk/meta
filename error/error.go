package error

import pb "meta/msg"

type MSGTYPE_NOT_FOUND struct{}

func (e MSGTYPE_NOT_FOUND) Error() string {
	return "MSGTYPE_NOT_FOUND"
}

type PROCESS_NOT_FOUND struct{}

func (e PROCESS_NOT_FOUND) Error() string {
	return "PROCESS_NOT_FOUND"
}

func GetErrorMsg(e error) (*pb.Msg, error) {
	return &pb.Msg{
		Type:    pb.ERR.Id,
		Content: []byte(e.Error()),
	}, e
}
