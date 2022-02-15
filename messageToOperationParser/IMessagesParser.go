package messageToOperationParser

type IMessagesParser interface {
	ParseMessage(message string)(IMapOperation, error)
}