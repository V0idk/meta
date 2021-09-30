package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type CommandExecutorConfig struct {
	Location string `json:"location"`
}

func GetCommandExecutorConfig(path string) *CommandExecutorConfig {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Printf("open %s error", jsonFile)
		return nil
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	commandExecutorConfig := CommandExecutorConfig{}
	err = json.Unmarshal(byteValue, &commandExecutorConfig)
	if err != nil {
		log.Printf("Unmarshal %s error", jsonFile)
		return nil
	}
	return &commandExecutorConfig
}
