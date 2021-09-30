package error

type MSGTYPE_NOT_FOUND struct{}

func (e MSGTYPE_NOT_FOUND) Error() string {
	return "MSGTYPE_NOT_FOUND"
}

type PROCESS_NOT_FOUND struct{}

func (e PROCESS_NOT_FOUND) Error() string {
	return "PROCESS_NOT_FOUND"
}
