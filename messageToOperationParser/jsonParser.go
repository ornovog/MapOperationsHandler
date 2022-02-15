package messageToOperationParser

import (
	"encoding/json"
)

type jsonParser struct {
}

type MessageBody struct {
	OperationType	string `json:"operation"`
	Key 		  	string `json:"key"`
	Value 		 	string `json:"value"`
}

func (_ jsonParser) ParseMessage(messageJson string) (IMapOperation, error){
	var messageBody MessageBody
	err := json.Unmarshal([]byte(messageJson), &messageBody)
	if err != nil{
		return nil, err
	}

	mapOp, err := MapOperation(messageBody)
	return mapOp, err

}

func MakeJsonParser() IMessagesParser {
	return jsonParser{}
}