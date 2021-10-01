package json

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func GetJsonBytes(path string) ([]byte, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Printf("open %s error", jsonFile)
		return nil, err
	}
	defer jsonFile.Close()
	return ioutil.ReadAll(jsonFile)
}

func ParseJson(path string, i *interface{}) error {
	byteValue, err := GetJsonBytes(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, &i)
	if err != nil {
		log.Printf("Unmarshal %s error", path)
		return err
	}
	return nil
}
