package main

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	pb "meta/msg"
	. "meta/processor/manager"
	. "meta/processor/manager/config"
	manager_msg "meta/processor/manager/msg"
	"net"
	"os"
)

type server struct {
	pb.UnimplementedMsgServiceServer
}

var m Manager

func (s *server) Dispatch(ctx context.Context, in *pb.Msg) (*pb.Msg, error) {
	log.Printf("Manage receive: %s", in)
	return m.Dispatch(in)
}

var managerConfig *ManagerConfig

func loadConfig() {
	managerConfig = GetManagerConfig(os.Args[1])
	if managerConfig == nil {
		log.Fatalf("failed to GetServerConfig")
	}
	m = Manager{
		Entry: manager_msg.Entry{
			Id:       uuid.New().String(),
			Location: managerConfig.Location,
		},
		Cache:         make(map[string]*EntryManager),
		ManagerConfig: *managerConfig,
	}
}

//GOMAXPROCS
func main() {
	loadConfig()
	lis, err := net.Listen("tcp", managerConfig.Location)
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
