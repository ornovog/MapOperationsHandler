package messageConsumerAndPublisher

import (
	"MapServer/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"strconv"
	"sync"
)

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
		QueueUrl:            aws.String(s.sqsConfig.QueueUrl),
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
			if !didMessageConsumeBefore(msg) {
				select {
				case s.msgChan <- msg.Body:
					msgsToDelete = append(msgsToDelete, msg)
				default:
				}
			}else {
				msgsToDelete = append(msgsToDelete, msg)
			}
		}

		s.deleteMessages(msgsToDelete)
	}
}

func didMessageConsumeBefore(msg *sqs.Message) bool{
	msgReceiveCountStr := msg.Attributes[sqs.MessageSystemAttributeNameApproximateReceiveCount]
	msgReceiveCount, _ := strconv.Atoi(*msgReceiveCountStr)
	return msgReceiveCount > 1
}

func (s *sqsConsumerAndPublisher) deleteMessages(messages []*sqs.Message){
	deleteMessageInput := &sqs.DeleteMessageInput{
		QueueUrl:      &s.sqsConfig.QueueUrl,
	}

	for _, msg := range messages {
		deleteMessageInput.ReceiptHandle = msg.ReceiptHandle
		_, err := s.sqsClient.DeleteMessage(deleteMessageInput)
		if err != nil {
			log.Printf("error when deleting - %s\n", err.Error())
		}
	}
}

//func (s *sqsConsumerAndPublisher) deleteMessages(messages []*sqs.Message){
//	deleteMessageInput := &sqs.DeleteMessageInput{
//		QueueUrl:      &s.sqsConfig.QueueUrl,
//	}
//
//	var wg sync.WaitGroup
//	for _, msg := range messages {
//		wg.Add(1)
//		deleteMessageInput.ReceiptHandle = msg.ReceiptHandle
//		go func(){
//			_, err := s.sqsClient.DeleteMessage(deleteMessageInput)
//			if err != nil {
//				log.Printf("error when deleting - %s\n", err.Error())
//			}
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//
// }


func MakeSqsConsumerAndPublisher(msgChan chan *string, sqsConfig config.SqsConfig) IMessagesConsumerAndPublisher{
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