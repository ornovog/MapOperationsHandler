package client

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func ReadFromCmdLine() {
	for true {
		line, err := getOperationFromCmdLine()
		if err != nil {
			log.Println(err.Error())
		} else {
			Execute(line)
		}
	}
}

func getOperationFromCmdLine() (string, error) {
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	line = strings.Split(line, "\n")[0]
	return line, err
}

func WriteRandomCommendsWithSleepTime(sleepingTime time.Duration){
	for ;true;{
		cmd := getRandomOperation()
		go Execute(cmd)
		time.Sleep(sleepingTime)
	}
}

func getRandomOperation() string{
	key := rand.Int()%1000
	val := rand.Int()%1000
	operation := rand.Int()%4
	switch operation {
	case 0:
		return fmt.Sprintf("add %d, %d", key, val)
	case 1:
		return fmt.Sprintf("get %d", key)
	case 2:
		return  fmt.Sprintf("remove %d", key)
	default:
		return fmt.Sprintf("getAll")
	}
}
