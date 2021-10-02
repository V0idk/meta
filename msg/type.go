package msg

type MsgType struct {
	Id string
}

var HEARTBEAT = MsgType{Id: "heartbeat"}
var REGISTER = MsgType{Id: "register"}
var BATCH = MsgType{Id: "batch"}
var COMMAND = MsgType{Id: "command"}
var OK = MsgType{Id: "ok"}
var ERR = MsgType{Id: "err"}
