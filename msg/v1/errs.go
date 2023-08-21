package v1

import (
	"errors"
)

var (
	errConfigNull         = errors.New("config cannot be null")
	errConfigServersEmpty = errors.New("servers cannot be empty")
	errNatsConnect        = errors.New("nats cannot connect")
	errProducerStart      = errors.New("start producer failed")
	errAdminClientStart   = errors.New("start admin client failed")
	errUnmarshal          = errors.New("unmarshal failed")
	errPublish            = errors.New("publish message failed")
)
