package v1

import (
	"errors"
)

var (
	errConfigNull          = errors.New("config cannot be null")
	errConfigFilePathEmpty = errors.New("filePath cannot be empty")
)
