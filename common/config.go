package common

import (
	"encoding/json"
	"os"
)

func SaveConfig(file string, config interface{}) {
	configFile, _ := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC, 0755) //Need to look into FileModes and general UNIX file permissions.
	defer configFile.Close()
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(config)
	if err != nil {
		panic(err)
	}
}

func LoadConfig(file string, config interface{}) error {
	configFile, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0755) //Need to look into FileModes and general UNIX file permissions.
	if err != nil {
		return err
	}
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	return decoder.Decode(config)
}
