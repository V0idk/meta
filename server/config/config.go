package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type MsgTypeConfig struct {
	Type string `json:"type"` //消息类型映射是server的责任，处理器们并不关心映射关系。因此放在这
	Rpc  string `json:"rpc"`
}

type RpcConfig struct {
	Name string `json:"name"`
	Type string `json:"type"`
	//https://stackoverflow.com/questions/58073214/is-there-a-way-to-dynamically-unmarshal-json-base-on-content
	//动态加载，根据rpc type决定
	Param   json.RawMessage `json:"Param"` //取决于Rpc的类型参数.
	Command string          `json:"command"`
	Args    []string        `json:"args"`
}

type ServerConfig struct {
	Location string          `json:"Location"`
	Msgtype  []MsgTypeConfig `json:"msgtype"`
	Rpc      []RpcConfig     `json:"rpc"`
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
