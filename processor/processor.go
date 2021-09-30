package processor

import (
	pb "meta/msg"
)

//====================================处理器类型，处理器通信方式==========================================
type Processor interface {
	Send(in *pb.Msg) (*pb.Msg, error)
}
