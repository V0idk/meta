package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	self_error "meta/error"
	pb "meta/msg"
	. "meta/processor"
	. "meta/processor/grpc"
	"meta/server/config"
	"net"
	"os"
)

//====================================配置文件加载==========================================
var serverConfig *config.ServerConfig
var msgTypeMap = make(map[string]config.MsgTypeConfig)
var processConfigMap = make(map[string]config.ProcessConfig)
var processMap = make(map[string]Processor)

func loadConfig() {
	serverConfig = config.GetServerConfig(os.Args[1])
	if serverConfig == nil {
		log.Fatalf("failed to GetServerConfig")
	}
	for _, item := range serverConfig.Msgtype {
		msgTypeMap[item.Type] = item
	}
	for _, item := range serverConfig.Process {
		processConfigMap[item.Type] = item
	}
	for _, v := range processConfigMap {
		if v.Type == "grpc" {
			processMap[v.Name] = &GrpcProcessor{
				Location: processConfigMap["grpc"].Args,
			}
		} else {
			log.Printf("Unsupport %s", v.Type)
		}
	}
}

//GOMAXPROCS
//=============================消息分发主进程======================================
type server struct {
	pb.UnimplementedMsgServiceServer
}

// 这个函数是并发的，因此这个函数可以阻塞。这个函数也必须阻塞，因为需要获取结果返回给对应的调用者。
func (s *server) Dispatch(ctx context.Context, in *pb.Msg) (*pb.Msg, error) {
	log.Printf("Server start to dispatch %s", in)
	if _, ok := msgTypeMap[in.Type]; !ok {
		log.Printf("Failed to find msgtype: %s", in.Type)
		return nil, &self_error.MSGTYPE_NOT_FOUND{}
	}

	if _, ok := processMap[msgTypeMap[in.Type].Process]; !ok {
		log.Printf("Failed to find process: %s", in.Type)
		return nil, &self_error.PROCESS_NOT_FOUND{}
	}
	process := processMap[msgTypeMap[in.Type].Process]
	msg, err := process.Send(in)
	if err != nil {
		log.Printf("Server fail to dispatch %s", in)
		return nil, nil
	}
	log.Printf("Server succues to dispatch %s", in)
	return msg, err
}

func main() {
	//cmd := exec.Command("manager/main/main.exe", "50001")
	//cmd.Start()
	loadConfig()
	lis, err := net.Listen("tcp", serverConfig.Location)
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
