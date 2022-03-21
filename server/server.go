package server

import (
	"MapServer/server/config"
	logger "MapServer/server/logger"
	msgConsAndPub "MapServer/server/messageConsumerAndPublisher"
	"MapServer/server/messageToOperationParser"
	"MapServer/server/operatinWorker"
	storage "MapServer/server/storage"
	"sync"
)

var (
	once sync.Once
)

func Serve() {
	once.Do(func() {
		configs, _ := config.GetConfigs()
		serverConfig := configs.ServerConfig
		msgChan := make(chan *string, serverConfig.ChannelSize)
		parser := messageToOperationParser.MakeJsonParser()
		storage := storage.MakeNavigationMap()
		logger := logger.MakeFileLogger(configs.LogsConfig.LogFilePath)
		runMapWorkers(serverConfig.NumOfMapWorkers, msgChan, parser, storage, logger)
		runMessageQueueConsumers(serverConfig.NumOfQueueConsumers, msgChan, configs.SqsConfig)
		select {}
	})
}

func runMapWorkers(numOfWorkers int, msgChan chan *string, parser messageToOperationParser.IMessagesParser,
	storage storage.IStorage, logger logger.ILogger) {
	for i := 0; i < numOfWorkers; i++ {
		worker := operatinWorker.MakeMapOperationsWorker(msgChan, parser, storage, logger)
		go worker.ReadMessagesAndExecuteMapOperations()
	}
}

func runMessageQueueConsumers(numOfConsumers int, msgChan chan *string, sqsConfig config.SqsConfig) {
	for i := 0; i < numOfConsumers; i++ {
		consumer := msgConsAndPub.MakeSqsConsumerAndPublisher(msgChan, sqsConfig)
		go consumer.StartConsuming()
	}
}
