package config

import (
	"encoding/json"
	"os"
	"path"
)

type Configuration struct {
	ServerConfig ServerConfig
	SqsConfig    SqsConfig
	LogsConfig   LogsConfig
}

type ServerConfig struct{
	NumOfMapWorkers     int
	NumOfQueueConsumers	int
}

type SqsConfig struct{
	QueueUrl			string
	MaxNumberOfMessages int64
	WaitTimeSeconds		int64
	VisibilityTimeout	int64
}

type LogsConfig struct{
	LogFilePath	string
}

func GetConfigs() (Configuration, error){
	config := Configuration{}
	wd, err := os.Getwd()
	if err != nil{
		return config, err
	}

	configFile, err := os.Open(path.Join(wd, "config/config.json"))
	if err != nil{
		return config, err
	}

	defer configFile.Close()
	decoder := json.NewDecoder(configFile)

	err = decoder.Decode(&config)
	return config, err
}