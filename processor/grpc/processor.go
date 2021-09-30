package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "meta/msg"
	"time"
)

type GrpcProcessor struct {
	Location string
}

func (s *GrpcProcessor) Send(in *pb.Msg) (*pb.Msg, error) {
	conn, err := grpc.Dial(s.Location, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMsgServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg, err := c.Dispatch(ctx, in)
	return msg, err
}
