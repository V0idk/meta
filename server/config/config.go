package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type MsgTypeConfig struct {
	Type    string `json:"type"`
	Process string `json:"process"`
}

type ProcessConfig struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Command string `json:"command"`
	Args    string `json:"args"`
}

type ServerConfig struct {
	Location string          `json:"Location"`
	Msgtype  []MsgTypeConfig `json:"msgtype"`
	Process  []ProcessConfig `json:"process"`
}

func GetServerConfig(path string) *ServerConfig {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Printf("open %s error", jsonFile)
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	serverConfig := ServerConfig{}
	err = json.Unmarshal(byteValue, &serverConfig)
	if err != nil {
		log.Printf("Unmarshal %s error", jsonFile)
		return nil
	}
	return &serverConfig
}
