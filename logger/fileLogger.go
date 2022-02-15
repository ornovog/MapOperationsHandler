package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type fileLogger struct {
	counter  int
	mu 		 sync.Mutex
	filePath string
}

func (l *fileLogger) WriteToLog(msg string)error{
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	timeStamp := time.Now().Format(time.Stamp)
	str := fmt.Sprintf("line:%d, time:%s, - %s \n", l.counter, timeStamp, msg)
	l.counter++
	_, err = f.Write([]byte(str))
	return err
}

func MakeFileLogger(filePath string) ILogger {
	l := fileLogger{filePath: filePath}
	return &l
}