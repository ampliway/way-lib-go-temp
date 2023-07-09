package v1

import (
	"errors"
	"fmt"
)

var (
	errConfigLen              = fmt.Errorf("configs must be equal or less than %v", configMaxLen)
	errConfigAlreadyExist     = errors.New("config already exist")
	errConfigNil              = errors.New("config cannot be nil")
	errConfigNameEmpty        = errors.New("config name cannot be empty")
	errConfigNameLen          = fmt.Errorf("config name cannot be more than %v", configNameMaxLen)
	errConfigDescriptionEmpty = errors.New("config description cannot be empty")
	errConfigDescriptionLen   = fmt.Errorf("config description cannot be more than %v", configDescriptionMaxLen)
	errConfigExecuteNil       = errors.New("config execute cannot be nil")
	errEmptyArguments         = errors.New("with empty arguments")
	errUnknown                = errors.New("unknown command")
	errExecutionFailed        = errors.New("execution failed")
)
