package config

import (
	"encoding/json"
	"log"
	. "meta/processor/manager/msg"
	. "meta/utils/json"
	"time"
)

type ManagerConfig struct {
	Entry         Entry         `json:"entry"`
	Location      string        `json:"location"`
	HeartbeatTime time.Duration `json:"heartbeat_time"`
}

func GetManagerConfig(path string) *ManagerConfig {
	managerConfig := ManagerConfig{}
	byteValue, err := GetJsonBytes(path)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(byteValue, &managerConfig)
	if err != nil {
		log.Printf("Unmarshal %s error", path)
		return nil
	}
	return &managerConfig
}
