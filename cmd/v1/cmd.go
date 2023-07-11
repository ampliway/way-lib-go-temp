package v1

import (
	"fmt"
	"os"
	"strings"

	"github.com/ampliway/way-lib-go/app"
	appV1 "github.com/ampliway/way-lib-go/app/v1"
	"github.com/ampliway/way-lib-go/cmd"
)

const (
	configMaxLen            = 10
	configNameMaxLen        = 20
	configDescriptionMaxLen = 300
)

var _ cmd.V1[any] = (*Cmd[any])(nil)

type Cmd[T any] struct {
	configs []*cmd.Config[T]
}

func New[T any]() *Cmd[T] {
	cmd := &Cmd[T]{
		configs: []*cmd.Config[T]{},
	}

	cmd.addReservedCommands()

	return cmd
}

func (c *Cmd[T]) Add(config *cmd.Config[T]) error {
	if len(c.configs)+1 > configMaxLen {
		return fmt.Errorf("%s: %w", cmd.MODULE_NAME, errConfigLen)
	}

	if err := configIsValid(config); err != nil {
		return err
	}

	for _, checkConfig := range c.configs {
		if config.Name == checkConfig.Name {
			return fmt.Errorf("%s: %w: %s", cmd.MODULE_NAME, errConfigAlreadyExist, config.Name)
		}
	}

	c.configs = append(c.configs, config)

	return nil
}

func (c *Cmd[T]) Run(arguments ...string) {
	appModule, err := appV1.New[T]()
	if err != nil {
		panic(fmt.Errorf("%s: %w", cmd.MODULE_NAME, err))
	}

	customArguments := strings.Join(arguments, " ")
	if customArguments == "" {
		customArguments = strings.Join(os.Args[1:], " ")
	}

	args := replaceSpaceSplit(customArguments)
	if len(args) == 0 || args[0] == "" {
		panic(fmt.Errorf("%s: %w", cmd.MODULE_NAME, errEmptyArguments))
	}

	match := c.findConfig(args...)

	if match == nil {
		panic(fmt.Errorf("%s: %w: %v", cmd.MODULE_NAME, errUnknown, args))
	}

	err = match.Execute(appModule)
	if err != nil {
		panic(fmt.Errorf("%s: %w: %+v", cmd.MODULE_NAME, errExecutionFailed, match))
	}
}

func configIsValid[T any](config *cmd.Config[T]) error {
	if config == nil {
		return fmt.Errorf("%s: %w", cmd.MODULE_NAME, errConfigNil)
	}

	config.Name = strings.TrimSpace(config.Name)
	if config.Name == "" {
		return fmt.Errorf("%s: %w", cmd.MODULE_NAME, errConfigNameEmpty)
	}

	if len(config.Name) > configNameMaxLen {
		return fmt.Errorf("%s: %w", cmd.MODULE_NAME, errConfigNameLen)
	}

	config.Args = replaceSpaceSplit(config.Name)

	config.Description = strings.TrimSpace(config.Description)
	if config.Description == "" {
		return fmt.Errorf("%s: %w", cmd.MODULE_NAME, errConfigDescriptionEmpty)
	}

	if len(config.Description) > configDescriptionMaxLen {
		return fmt.Errorf("%s: %w", cmd.MODULE_NAME, errConfigDescriptionLen)
	}

	if config.Execute == nil {
		return fmt.Errorf("%s: %w", cmd.MODULE_NAME, errConfigExecuteNil)
	}

	return nil
}

func replaceSpaceSplit(input string) []string {
	input = strings.TrimSpace(input)
	oldString := "  "
	newString := " "

	for strings.Contains(input, oldString) {
		input = strings.ReplaceAll(input, oldString, newString)
	}

	if input == "" {
		return []string{}
	}

	return strings.Split(input, newString)
}

func (c *Cmd[T]) findConfig(args ...string) *cmd.Config[T] {
	var match *cmd.Config[T]

	for _, config := range c.configs {
		if len(args) != len(config.Args) {
			continue
		}

		isEqual := true

		for i := 0; i < len(args); i++ {
			if config.Args[i] != args[i] {
				isEqual = false

				break
			}
		}

		if isEqual {
			match = config

			break
		}
	}

	return match
}

func (c *Cmd[T]) addReservedCommands() {
	c.configs = append(c.configs, &cmd.Config[T]{
		Name:        "commands",
		Description: "List all commands",
		Execute: func(app app.V1[T]) error {
			for _, command := range c.configs {
				fmt.Println(command.Name)
			}

			return nil
		},
		Args: []string{},
	})
}
