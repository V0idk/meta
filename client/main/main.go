package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	pb "meta/msg"
	. "meta/processor/command_executor/msg"
	. "meta/processor/manager/msg"
	"strconv"
	"sync"
)

func dial(type_param string, content []byte) {
	address := fmt.Sprintf("127.0.0.1:50000")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		wg.Done()
		return
	}
	defer conn.Close()
	c := pb.NewMsgServiceClient(conn)
	msg, err := c.Dispatch(context.Background(), &pb.Msg{
		Type:    type_param,
		Content: content,
	})
	log.Printf("dial result: %s, err: %s", msg, err)
}

var wg sync.WaitGroup

func query(i int) {
	dial(strconv.Itoa(i), []byte(""))
	wg.Done()
}

func testQuery() {
	for i := 1; i <= 300; i++ {
		wg.Add(1)
		go query(i)
	}
	wg.Wait()
}

//==============================
func testManager() {
	content := &HeartbeatContent{}
	content.Entry.Id = uuid.New().String()
	content.Entry.Location = "127.0.0.1:50001"
	result, err := json.Marshal(content)
	if err != nil {
		log.Printf("could not marshal: %s", content)
		return
	}
	dial(pb.HEARTBEAT.Id, result)
}

func testCommand() {
	content := &CommandContent{}
	content.Command = "cmd"
	content.Args = append(content.Args, "/C")
	content.Args = append(content.Args, "dir")
	result, err := json.Marshal(content)
	if err != nil {
		log.Printf("could not marshal: %s", content)
		return
	}
	dial(pb.COMMAND.Id, result)
}

func main() {
	testCommand()
}
