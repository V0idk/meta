package msg

import pb "meta/msg"

type Entry struct {
	Id       string `json:"id"`
	Location string `json:"location"`
}
type HeartbeatContent struct {
	Entry Entry `json:"entry"`
}

type RegisterContent struct {
	Entry Entry `json:"entry"`
}

type BatchType struct {
	Id string
}

var LIST = BatchType{Id: "list"}
var ALL = BatchType{Id: "ALL"}

type PairResult struct {
	Entry Entry   `json:"entry"`
	Msg   *pb.Msg `json:"msg"`
	Err   error   `json:"err"`
}

type BatchContent struct {
	Type   string  `json:"type"`
	Entrys []Entry `json:"entrys"`
	Msg    *pb.Msg `json:"msg"`
}

type BatchResult struct {
	PairResults []PairResult `json:"pairResults"`
}
