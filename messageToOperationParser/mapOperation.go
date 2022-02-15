package messageToOperationParser

type mapOperation struct {
	operationType OperationType
	key           string
	value         string
}

func (m mapOperation) OperationType() OperationType {
	return m.operationType
}

func (m mapOperation) Key() string{
	return m.key
}

func (m mapOperation) Value() string{
	return m.value
}

func MapOperation(messageBody MessageBody) (IMapOperation, error){
	m := mapOperation{}
	opType, err := stringToOperation(messageBody.OperationType)
	if err != nil{
		return nil, err
	}

	m.operationType = opType
	m.key = messageBody.Key
	m.value = messageBody.Value
	return m, nil
}
