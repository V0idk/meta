package config

import (
	"encoding/json"
	"log"
	. "meta/utils/json"
)

type CommandExecutorConfig struct {
	Location string `json:"location"`
}

func GetCommandExecutorConfig(path string) *CommandExecutorConfig {
	commandExecutorConfig := CommandExecutorConfig{}
	byteValue, err := GetJsonBytes(path)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(byteValue, &commandExecutorConfig)
	if err != nil {
		log.Printf("Unmarshal %s error", path)
		return nil
	}
	return &commandExecutorConfig
}
