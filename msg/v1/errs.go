package v1

import (
	"errors"
)

var (
	errConfigNull         = errors.New("config cannot be null")
	errConfigServersEmpty = errors.New("servers cannot be empty")
	errNATSConnect        = errors.New("nats cannot connect")
	errJSConnect          = errors.New("jetstream cannot connect")
	errUnmarshal          = errors.New("unmarshal failed")
	errSubPrefix          = errors.New("invalid prefix message")
	errPublish            = errors.New("publish message failed")
)
