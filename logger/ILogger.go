package logger

type ILogger interface {
	WriteToLog(str string) error
}
