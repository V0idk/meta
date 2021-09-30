package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	pb "meta/msg"
	manager_msg "meta/processor/manager/msg"
	"strconv"
	"sync"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg, err := c.Dispatch(ctx, &pb.Msg{
		Type:    type_param,
		Content: content,
	})
	if err != nil {
		log.Printf("could not dial: %v", err)
		return
	}
	log.Printf("dial result: %s", msg)
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
	content := &manager_msg.HeartbeatContent{}
	content.Entry.Id = uuid.New().String()
	content.Entry.Location = "127.0.0.1:50001"
	result, err := json.Marshal(content)
	if err != nil {
		log.Printf("could not marshal: %s", content)
		return
	}
	dial(pb.HEARTBEAT.Id, result)
}

func main() {
	testManager()
}
