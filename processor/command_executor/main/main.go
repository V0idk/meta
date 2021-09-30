package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "meta/msg"
	. "meta/processor/command_executor"
	. "meta/processor/command_executor/config"
	"net"
	"os"
)

type server struct {
	pb.UnimplementedMsgServiceServer
}

var c CommandExecutor

func (s *server) Dispatch(ctx context.Context, in *pb.Msg) (*pb.Msg, error) {
	log.Printf("Manage receive: %s", in)
	return c.Dispatch(in)
}

var commandExecutorConfig *CommandExecutorConfig

func loadConfig() {
	commandExecutorConfig = GetCommandExecutorConfig(os.Args[1])
	if commandExecutorConfig == nil {
		log.Fatalf("failed to GetServerConfig")
	}
	c = CommandExecutor{}
}

//GOMAXPROCS
func main() {
	loadConfig()
	lis, err := net.Listen("tcp", commandExecutorConfig.Location)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMsgServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
