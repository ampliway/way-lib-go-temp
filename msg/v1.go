package msg

import "github.com/ampliway/way-lib-go/ctx"

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
	Publish(ctx ctx.V1, m interface{}) error
	Shutdown()
}

type SubscriberV1[T any] interface {
	Publish(ctx ctx.V1, m interface{}) error
	Subscribe(queueGroup string, execution func(msg *Message[T]) bool) error
	Shutdown()
}
