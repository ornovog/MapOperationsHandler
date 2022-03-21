package messageToOperationParser

import "fmt"

type OperationType int64
const (
	Add OperationType = iota
	Remove
	Get
	GetAll
)

type IMapOperation interface {
	OperationType() OperationType
	Key() string
	Value() string
}

func stringToOperation(operation string)(OperationType, error){
	switch operation{
	case "add":
		return Add, nil
	case "remove":
		return Remove, nil
	case "get":
		return Get, nil
	case "getAll":
		return GetAll, nil
	default:
		return 0, fmt.Errorf("operation %s doesn't exist", operation)
	}
}