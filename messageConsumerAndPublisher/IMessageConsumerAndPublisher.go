package messageConsumerAndPublisher

type IMessagePublisher interface {
	GetMsgPublishChannel() chan *string
}

type IMessagesConsumerAndPublisher interface {
	GetMsgPublishChannel() chan *string
	StartConsuming()
}


