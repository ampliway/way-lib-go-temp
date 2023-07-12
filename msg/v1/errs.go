package v1

import (
	"errors"
)

var (
	errConfigNull         = errors.New("config cannot be null")
	errConfigServersEmpty = errors.New("servers cannot be empty")
	errKafkaConnect       = errors.New("kafka cannot connect")
	errProducerStart      = errors.New("start producer failed")
	errAdminClientStart   = errors.New("start admin client failed")
	errTopicCreate        = errors.New("create topic failed")
	errUnmarshal          = errors.New("unmarshal failed")
	errSubPrefix          = errors.New("invalid prefix message")
	errPublish            = errors.New("publish message failed")
)
