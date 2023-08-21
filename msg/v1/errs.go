package v1

import (
	"errors"
)

var (
	errConfigNull         = errors.New("config cannot be null")
	errConfigServersEmpty = errors.New("servers cannot be empty")
	errNatsConnect        = errors.New("nats cannot connect")
	errUnmarshal          = errors.New("unmarshal failed")
	errPublish            = errors.New("publish message failed")
)
