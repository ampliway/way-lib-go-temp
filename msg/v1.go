package msg

const MODULE_NAME = "msg"

type MsgV1 interface {
	Publish(m interface{}) error
	PublishT(topicName string, m interface{}) error
	Subscribe(m interface{}, queueGroup string, exec func(data []byte) bool) error
	CreateTopicIfNotExist(topicName string) error
	Shutdown()
}
