package manager

import (
	"encoding/json"
	"log"
	self_error "meta/error"
	pb "meta/msg"
	. "meta/processor/manager/config"
	manager_msg "meta/processor/manager/msg"
	. "meta/rpc"
	. "meta/rpc/grpc"
	"time"
)

//同一个文件夹下的包只能依赖更内层或同级的包，而不能依赖外层包，比如A/B/in.go依赖A/out.go
//main除外
type EntryManager struct {
	Entry       manager_msg.Entry
	Rpc         Rpc
	SuccessTime time.Time
}

func (e *EntryManager) DialHeartbeat(src manager_msg.Entry) (*pb.Msg, error) {
	log.Printf("Start to DialHeartbeat: %s", e.Entry)
	content := &manager_msg.HeartbeatContent{}
	content.Entry = src
	result, err := json.Marshal(content)
	if err != nil {
		log.Printf("could not marshal: %s", content)
		return nil, err
	}
	msg, err := e.Rpc.Send(&pb.Msg{
		Type:    pb.HEARTBEAT.Id,
		Content: result,
	})
	if err != nil {
		log.Printf("Failed to DialHeartbeat: %s", err)
		return nil, err
	}
	log.Printf("Success to DialHeartbeat: %s", e.Entry)
	return msg, err
}

func (e *EntryManager) HeartbeatSuccess() {
	e.SuccessTime = time.Now()
}

type Manager struct {
	Entry         manager_msg.Entry
	Cache         map[string]*EntryManager
	ManagerConfig ManagerConfig
}

//心跳函数，也是注册函数。注册的发起由外部程序管理控制，比如client。
func (m *Manager) Heartbeat(in *pb.Msg) (*pb.Msg, error) {
	var content manager_msg.HeartbeatContent
	err := json.Unmarshal(in.Content, &content)
	if err != nil {
		log.Printf("Failed to load")
		return nil, err
	}
	if _, ok := m.Cache[content.Entry.Id]; !ok {
		// 不存在则增加entry，且启动后台心跳检测
		tmp := &EntryManager{
			Entry:       content.Entry,
			Rpc:         &Grpc{content.Entry.Location},
			SuccessTime: time.Now(),
		}
		m.Cache[content.Entry.Id] = tmp
		go m.DialHeartbeatBackground(m.Cache[content.Entry.Id])
		log.Printf("Entry %s heartbeat success", content.Entry)
	} else {
		//存在则刷新状态
		log.Printf("Entry %s heartbeat success", content.Entry)
		//https://stackoverflow.com/questions/32751537/why-do-i-get-a-cannot-assign-error-when-setting-value-to-a-struct-as-a-value-i
		m.Cache[content.Entry.Id].HeartbeatSuccess()
	}
	return &pb.Msg{Type: pb.OK.Id}, nil
}

//为什么DialHeartbeatBackground不放在EntryManager? 因为Heartbeat函数修改了EntryManager状态。所以我们让EntryManager保持纯粹性，只提供属性和方法。
func (m *Manager) DialHeartbeatBackground(e *EntryManager) {
	for {
		_, err := e.DialHeartbeat(m.Entry)
		if err == nil {
			e.HeartbeatSuccess()
		}
		time.Sleep(m.ManagerConfig.HeartbeatTime * time.Millisecond)
	}
}

func (m *Manager) Dispatch(in *pb.Msg) (*pb.Msg, error) {
	switch in.Type {
	case pb.HEARTBEAT.Id:
		return m.Heartbeat(in)
	}
	return nil, &self_error.MSGTYPE_NOT_FOUND{}
}
