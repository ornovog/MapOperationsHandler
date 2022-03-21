package client

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"strings"
)

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	sqsClient = sqs.New(sess)
	qURL = os.Getenv("qURL")
)

type Operation struct {
	OperationType	string `json:"operation"`
	Key 		  	string `json:"key"`
	Value 		 	string `json:"value"`
}

func Execute (str string){
	operation, err := stringToOperation(str)
	if err != nil{
		log.Print(err.Error())
		return
	}

	messageBody, err := operationToJsonStr(operation)
	if err != nil{
		log.Print(err.Error())
		return
	}

	receiveMessageInput := sqs.SendMessageInput{
		QueueUrl:    &qURL,
		MessageBody: &messageBody,
	}

	_, err = sqsClient.SendMessage(&receiveMessageInput)
	if err != nil{
		log.Print(err.Error())
		return
	}
}

func stringToOperation(str string) (Operation, error) {
	cmd := strings.Split(str, " ")
	invalidCmdErr := fmt.Errorf("invalid command - %s", str)
	if len(cmd) == 0 {
		return Operation{}, invalidCmdErr
	}

	switch cmd[0] {
	case "add":
		if len(cmd) != 3{
			return Operation{}, invalidCmdErr
		}
		return Operation{OperationType: cmd[0], Key: cmd[1], Value: cmd[2]}, nil
	case "remove", "get":
		if len(cmd) != 2{
			return Operation{}, invalidCmdErr
		}
		return Operation{OperationType: cmd[0], Key: cmd[1]}, nil
	case "getAll":
		if len(cmd) != 1{
			return Operation{}, invalidCmdErr
		}
		return Operation{OperationType: cmd[0]}, nil
	default:
		return Operation{}, invalidCmdErr
	}
}

func operationToJsonStr(operation Operation) (string, error) {
	bytes, err := json.Marshal(operation)
	if err != nil {
		return "", err
	}
	messageBody := string(bytes)
	return messageBody, nil
}
