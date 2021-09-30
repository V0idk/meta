package msg

type MsgType struct {
	Id string
}

var HEARTBEAT = MsgType{Id: "heartbeat"}
var COMMAND = MsgType{Id: "command"}
var OK = MsgType{Id: "ok"}
