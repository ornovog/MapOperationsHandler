package messageConsumerAndPublisher

import (
	"MapServer/server/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"strconv"
	"sync"
)

var qURL = os.Getenv("qURL")

type sqsConsumerAndPublisher struct {
	sqsClient *sqs.SQS
	sqsConfig config.SqsConfig
	msgChan   chan *string
	once 	  sync.Once
}

func (s *sqsConsumerAndPublisher) GetMsgPublishChannel() chan *string {
	return s.msgChan
}

func (s *sqsConsumerAndPublisher) StartConsuming(){
	s.once.Do(func() {
		s.consumeAndPublishMessages()
	})
}

func (s *sqsConsumerAndPublisher) consumeAndPublishMessages(){
	receiveMessageInput := sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(qURL),
		MaxNumberOfMessages: aws.Int64(s.sqsConfig.MaxNumberOfMessages),
		WaitTimeSeconds:     aws.Int64(s.sqsConfig.WaitTimeSeconds),
		AttributeNames: 	 aws.StringSlice([]string{sqs.MessageSystemAttributeNameApproximateReceiveCount}),
		VisibilityTimeout:   aws.Int64(s.sqsConfig.VisibilityTimeout),
	}

	for true {
		output, err := s.sqsClient.ReceiveMessage(&receiveMessageInput)
		if err != nil {
			log.Printf("error when consuming - %s\n", err.Error())
		}

		msgsToDelete := make([]*sqs.Message, 0)
		for _, msg := range output.Messages {
			select {
			case s.msgChan <- msg.Body:
				msgsToDelete = append(msgsToDelete, msg)
			default:
			}
		}

		if len(msgsToDelete) > 0{
			s.deleteMessages(msgsToDelete)
		}
	}
}

func messageConsumedBefore(msg *sqs.Message) bool{
	msgReceiveCountStr := msg.Attributes[sqs.MessageSystemAttributeNameApproximateReceiveCount]
	msgReceiveCount, _ := strconv.Atoi(*msgReceiveCountStr)
	return msgReceiveCount > 1
}

func (s *sqsConsumerAndPublisher) deleteMessages(messages []*sqs.Message){
	deleteMessagesList := messagesToDeleteMessages(messages)
	deleteMessageInput := &sqs.DeleteMessageBatchInput{
		Entries:  deleteMessagesList,
		QueueUrl: &s.sqsConfig.QueueUrl,
	}

	output, err := s.sqsClient.DeleteMessageBatch(deleteMessageInput)
	if err != nil {
		log.Printf("error when deleting - %s\n", err.Error())
		return
	}

	for _, failed := range output.Failed{
		log.Printf("error when deleting - %s\n", *failed.Message)
	}
}

func messagesToDeleteMessages(messages []*sqs.Message) []*sqs.DeleteMessageBatchRequestEntry {
	deleteMessagesList := make([]*sqs.DeleteMessageBatchRequestEntry, 0)
	for _, msg := range messages {
		deleteMessage := sqs.DeleteMessageBatchRequestEntry{
			Id:            msg.MessageId,
			ReceiptHandle: msg.ReceiptHandle,
		}
		deleteMessagesList = append(deleteMessagesList, &deleteMessage)
	}
	return deleteMessagesList
}


func MakeSqsConsumerAndPublisher(msgChan chan *string, sqsConfig config.SqsConfig) IMessagesConsumerAndPublisher {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sqsClient := sqs.New(sess)

	sqsCons := sqsConsumerAndPublisher{
		msgChan: msgChan,
		sqsClient: sqsClient,
		sqsConfig: sqsConfig,
	}
	return &sqsCons
}