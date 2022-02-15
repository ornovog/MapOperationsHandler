package operatinWorker

import (
	"MapServer/logger"
	parser "MapServer/messageToOperationParser"
	"MapServer/storage"
	"fmt"
	"log"
	"sync"
)

type mapOperationsWorker struct {
	messagesChannel   chan *string
	messagesParser    parser.IMessagesParser
	storage           storage.IStorage
	logger            logger.ILogger
	once 	  		  sync.Once
}

func (m *mapOperationsWorker) ReadMessagesAndExecuteMapOperations(){
	m.once.Do(func() {
		m.readAndExecuteMapOperations()
	})
}

func (m *mapOperationsWorker) readAndExecuteMapOperations(){
	for msg := range m.messagesChannel{
		m.handleMessage(msg)
	}
}

func (m *mapOperationsWorker) handleMessage(msg *string){
	operation, err := m.messagesParser.ParseMessage(*msg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	opResult, err := m.executeOperation(operation)
	if err != nil {
		log.Println(err.Error())
		return
	}

	m.logger.WriteToLog(opResult)
}

func (m *mapOperationsWorker) executeOperation(operation parser.IMapOperation)(string, error){
	switch operation.OperationType() {
		case parser.Add:
			key := operation.Key()
			val := operation.Value()
			err := m.storage.AddItem(key, val)
			if err != nil {
				return err.Error(), nil
				//return "", err
			}
			return fmt.Sprintf("added - {%s : %s}", key, val), nil
		case parser.Remove:
			key := operation.Key()
			m.storage.RemoveItem(key)
			return fmt.Sprintf("removed key - %s", key), nil
		case parser.Get:
			key := operation.Key()
			val, err := m.storage.GetItem(key)
			if err != nil {
				return err.Error(), nil
				//return "", err
			}
			return fmt.Sprintf("got key:value - {%s : %s}", key, val), nil
		case parser.GetAll:
			allItems := m.storage.GetAllItems()
			return fmt.Sprintf("all items - %s", allItems), nil
		default:
			return "", fmt.Errorf("invalid operation - %v", operation.OperationType())

	}
}

func MakeMapOperationsWorker(messagesChanel chan *string, messagesParser parser.IMessagesParser,
	storage storage.IStorage, logger logger.ILogger)IMapOperationsWorker{

	m := mapOperationsWorker{
		messagesChannel : messagesChanel,
		messagesParser: messagesParser,
		storage : storage,
		logger : logger,
	}
	return &m
}