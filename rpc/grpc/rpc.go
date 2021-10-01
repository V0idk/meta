package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "meta/msg"
)

type Grpc struct {
	Location string
}

func (s *Grpc) Send(in *pb.Msg) (*pb.Msg, error) {
	conn, err := grpc.Dial(s.Location, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMsgServiceClient(conn)
	msg, err := c.Dispatch(context.Background(), in)
	return msg, err
}
