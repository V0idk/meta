package manager

import (
	"encoding/json"
	"log"
	. "meta/error"
	pb "meta/msg"
	. "meta/processor/manager/config"
	. "meta/processor/manager/msg"
	. "meta/rpc"
	. "meta/rpc/grpc"
	"sync"
	"time"
)

//同一个文件夹下的包只能依赖更内层或同级的包，而不能依赖外层包，比如A/B/in.go依赖A/out.go
//main除外
type EntryManager struct {
	Entry       Entry
	Rpc         Rpc
	SuccessTime time.Time
}

func (e *EntryManager) DialHeartbeat(src Entry) (*pb.Msg, error) {
	log.Printf("Start to DialHeartbeat: %s", e.Entry)
	content := &HeartbeatContent{}
	content.Entry = src
	result, err := json.Marshal(content)
	if err != nil {
		log.Printf("could not marshal: %s", content)
		return GetErrorMsg(err)
	}
	msg, err := e.Rpc.Send(&pb.Msg{
		Type:    pb.HEARTBEAT.Id,
		Content: result,
	})
	if err != nil {
		log.Printf("Failed to DialHeartbeat: %s", err)
		return GetErrorMsg(err)
	}
	log.Printf("Success to DialHeartbeat: %s", e.Entry)
	return msg, err
}

func (e *EntryManager) HeartbeatSuccess() {
	e.SuccessTime = time.Now()
}

type Manager struct {
	Entry         Entry
	Cache         map[string]*EntryManager
	ManagerConfig ManagerConfig
}

//心跳函数。
func (m *Manager) Heartbeat(in *pb.Msg) (*pb.Msg, error) {
	var content HeartbeatContent
	err := json.Unmarshal(in.Content, &content)
	if err != nil {
		log.Printf("Failed to load")
		return GetErrorMsg(err)
	}
	if _, ok := m.Cache[content.Entry.Id]; ok {
		m.Cache[content.Entry.Id].HeartbeatSuccess()
	}
	//存在则刷新状态
	log.Printf("Entry from %s heartbeat success", content.Entry)
	//https://stackoverflow.com/questions/32751537/why-do-i-get-a-cannot-assign-error-when-setting-value-to-a-struct-as-a-value-i
	return &pb.Msg{Type: pb.OK.Id}, nil
}

//心跳和注册函数不合一，因为会导致管理者的心跳请求爆发式增长。注册的发起由外部程序管理控制，比如client
func (m *Manager) Register(in *pb.Msg) (*pb.Msg, error) {
	var content HeartbeatContent
	err := json.Unmarshal(in.Content, &content)
	if err != nil {
		log.Printf("Failed to load")
		return GetErrorMsg(err)
	}
	if len(content.Entry.Id) == 0 {
		return GetErrorMsg(PARAM_WRONG{Param: "Id"})
	}
	if len(content.Entry.Location) == 0 {
		return GetErrorMsg(PARAM_WRONG{Param: "Location"})
	}
	if _, ok := m.Cache[content.Entry.Id]; !ok {
		// 不存在则增加entry，且启动后台心跳检测
		tmp := &EntryManager{
			Entry:       content.Entry,
			Rpc:         &Grpc{Location: content.Entry.Location},
			SuccessTime: time.Now(),
		}
		m.Cache[content.Entry.Id] = tmp
		go m.DialHeartbeatBackground(m.Cache[content.Entry.Id])
		log.Printf("Entry %s heartbeat success", content.Entry)
	}
	return &pb.Msg{Type: pb.OK.Id}, nil
}

//为什么DialHeartbeatBackground不放在EntryManager? 因为Heartbeat函数修改了EntryManager状态。所以我们让EntryManager保持纯粹性，只提供属性和方法。
func (m *Manager) DialHeartbeatBackground(e *EntryManager) {
	for {
		msg, _ := e.DialHeartbeat(m.Entry)
		if msg.Type == pb.OK.Id {
			e.HeartbeatSuccess()
		}
		time.Sleep(m.ManagerConfig.HeartbeatTime * time.Millisecond)
	}
}

func (m *Manager) GetAllEntry() []Entry {
	ret := make([]Entry, 0)
	for _, e := range m.Cache {
		ret = append(ret, e.Entry)
	}
	return ret
}

func (m *Manager) BatchDial(entrys []Entry, in *pb.Msg) []PairResult {
	var wg sync.WaitGroup
	// https://stackoverflow.com/questions/13670818/pair-tuple-data-type-in-go
	result := make(chan PairResult, len(entrys))
	count := 0
	for _, e := range entrys {
		if len(e.Location) == 0 {
			continue
		}
		go func(which Entry) {
			defer wg.Done()
			rpc := &Grpc{Location: which.Location}
			msg, err := rpc.Send(in)
			result <- PairResult{
				Entry: which,
				Msg:   msg,
				Err:   err,
			}
		}(e)
		wg.Add(1)
		count += 1
	}
	//https://stackoverflow.com/questions/44152988/append-not-thread-safe
	wg.Wait()
	resultList := make([]PairResult, 0)
	for i := 0; i < count; i++ {
		resultList = append(resultList, <-result)
	}
	return resultList
}

func (m *Manager) Batch(in *pb.Msg) (*pb.Msg, error) {
	var content BatchContent
	err := json.Unmarshal(in.Content, &content)
	if err != nil {
		log.Printf("Failed to load")
		return GetErrorMsg(err)
	}
	var results []PairResult
	switch content.Type {
	case ALL.Id:
		results = m.BatchDial(m.GetAllEntry(), content.Msg)
		break
	case LIST.Id:
		results = m.BatchDial(content.Entrys, content.Msg)
		break
	}
	batchResult := BatchResult{
		PairResults: results,
	}
	resultBytes, err := json.Marshal(batchResult)
	if err != nil {
		log.Printf("could not marshal: %s", batchResult)
		return GetErrorMsg(err)
	}
	return &pb.Msg{
		Type:    pb.OK.Id,
		Content: resultBytes,
	}, nil
}

func (m *Manager) Dispatch(in *pb.Msg) (*pb.Msg, error) {
	switch in.Type {
	case pb.HEARTBEAT.Id:
		return m.Heartbeat(in)
	case pb.BATCH.Id:
		return m.Batch(in)
	case pb.REGISTER.Id:
		return m.Register(in)
	}
	return GetErrorMsg(MSGTYPE_NOT_FOUND{})
}
