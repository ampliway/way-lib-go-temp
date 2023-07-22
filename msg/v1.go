package msg

const (
	MODULE_NAME       = "msg"
	HEADER_X_MSG_ID   = "x-msg-id"
	HEADER_X_TRACE_ID = "x-trace-id"
)

type Message[T any] struct {
	MessageID string
	TraceID   string
	Timestamp int64
	Body      T
}

type ProducerV1 interface {
	Publish(key string, m interface{}) error
	PublishT(topicName, key string, m interface{}) error
	CreateTopicIfNotExist(topicName string, numPartitions int32, replicationFactor int16) error
	Shutdown()
}

type SubscriberV1[T any] interface {
	Publish(key string, m interface{}) error
	PublishT(topicName, key string, m interface{}) error
	Subscribe(queueGroup string, execution func(msg *Message[T]) bool) error
	SubscribeT(topicName, queueGroup string, execution func(msg *Message[T]) bool) error
	Shutdown()
}
