package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type ManagerConfig struct {
	Location      string        `json:"location"`
	HeartbeatTime time.Duration `json:"heartbeat_time"`
}

func GetManagerConfig(path string) *ManagerConfig {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Printf("open %s error", jsonFile)
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	managerConfig := ManagerConfig{}
	err = json.Unmarshal(byteValue, &managerConfig)
	if err != nil {
		log.Printf("Unmarshal %s error", jsonFile)
		return nil
	}
	return &managerConfig
}
