package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	pb "meta/msg"
	. "meta/processor/command_executor/msg"
	. "meta/processor/manager/msg"
	"strconv"
	"sync"
)

func dial(type_param string, content []byte) (*pb.Msg, error) {
	address := fmt.Sprintf("127.0.0.1:50000")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v", err)
		wg.Done()
		return nil, nil
	}
	defer conn.Close()
	c := pb.NewMsgServiceClient(conn)
	msg, err := c.Dispatch(context.Background(), &pb.Msg{
		Type:    type_param,
		Content: content,
	})
	log.Printf("dial result: %s, err: %s", msg, err)
	return msg, err
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
func testManager(id string, location string) {
	content := &HeartbeatContent{}
	content.Entry.Id = id
	content.Entry.Location = location
	result, err := json.Marshal(content)
	if err != nil {
		log.Printf("could not marshal: %s", content)
		return
	}
	dial(pb.REGISTER.Id, result)
}

func testCommand() {
	content := &CommandContent{}
	content.Command = "go"
	content.Args = append(content.Args, "run")
	content.Args = append(content.Args, "processor\\manager\\main\\main.go")
	content.Args = append(content.Args, "processor\\manager\\config\\manager_50012.json")
	result, err := json.Marshal(content)
	if err != nil {
		log.Printf("could not marshal: %s", content)
		return
	}
	msg, err := dial(pb.COMMAND.Id, result)
	if err == nil {
		cmdResult := CommandResult{}
		err := json.Unmarshal(msg.Content, &cmdResult)
		if err != nil {
			log.Printf("%s", err)
		} else {
			log.Printf("%s", cmdResult)
		}
	}
}

func testBatch() {
	testManager("50001", "127.0.0.1:50001")
	testManager("50002", "127.0.0.1:50002")

	content := &CommandContent{}
	content.Command = "cmd"
	content.Args = append(content.Args, "/C")
	content.Args = append(content.Args, "dir")
	result, err := json.Marshal(content)
	if err != nil {
		log.Printf("could not marshal: %s", content)
		return
	}
	msg := pb.Msg{
		Type:    pb.COMMAND.Id,
		Content: result,
	}

	batchContent := BatchContent{
		Type: ALL.Id,
		Entrys: []Entry{
			{Id: "50001"},
			{Id: "50002"},
		},
		Msg: &msg,
	}
	batchContentResult, err := json.Marshal(batchContent)
	if err != nil {
		log.Printf("could not marshal: %s", content)
		return
	}
	dial(pb.BATCH.Id, batchContentResult)
}

func testFor() {
	type a struct {
		a1 string
	}
	al := []a{
		{"000"}, {"111"},
	}
	for _, e := range al {
		go func() {
			log.Printf("content: %s", e)
		}()
	}
}

//序列化reflect.Type?
//https://stackoverflow.com/questions/43770692/marshal-unmarshal-reflect-type
func main() {
	//testBatch()
}
