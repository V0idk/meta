package msg

type MsgType struct {
	Id string
}

var HEARTBEAT = MsgType{Id: "heartbeat"}
var OK = MsgType{Id: "ok"}
