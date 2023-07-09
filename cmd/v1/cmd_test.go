package v1

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/ampliway/way-lib-go/app"
	"github.com/ampliway/way-lib-go/cmd"
	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	Field1 string
}

func TestNew(t *testing.T) {
	t.Parallel()

	adapter := New[testConfig]()
	assert.NotNil(t, adapter)
	assert.NotNil(t, adapter.configs)
}

func TestAdd(t *testing.T) {
	t.Parallel()

	tableTest := []struct {
		Scenario    string
		Adapter     *Cmd[testConfig]
		Config      *cmd.Config[testConfig]
		ExpectedErr error
	}{
		{
			Scenario:    "config_nil",
			Adapter:     New[testConfig](),
			Config:      nil,
			ExpectedErr: errConfigNil,
		},
		{
			Scenario: "config_name_empty",
			Adapter:  New[testConfig](),
			Config: &cmd.Config[testConfig]{
				Name: "",
			},
			ExpectedErr: errConfigNameEmpty,
		},
		{
			Scenario: "config name empty with spaces",
			Adapter:  New[testConfig](),
			Config: &cmd.Config[testConfig]{
				Name: "    ",
			},
			ExpectedErr: errConfigNameEmpty,
		},
		{
			Scenario: "length of config name longer than allowed",
			Adapter:  New[testConfig](),
			Config: &cmd.Config[testConfig]{
				Name: strings.Repeat("a", configNameMaxLen+1),
			},
			ExpectedErr: errConfigNameLen,
		},
		{
			Scenario: "config description empty",
			Adapter:  New[testConfig](),
			Config: &cmd.Config[testConfig]{
				Name:        strings.Repeat("a", configNameMaxLen),
				Description: "",
			},
			ExpectedErr: errConfigDescriptionEmpty,
		},
		{
			Scenario: "config description empty with white spaces",
			Adapter:  New[testConfig](),
			Config: &cmd.Config[testConfig]{
				Name:        strings.Repeat("a", configNameMaxLen),
				Description: "   ",
			},
			ExpectedErr: errConfigDescriptionEmpty,
		},
		{
			Scenario: "length of config description longer than allowed",
			Adapter:  New[testConfig](),
			Config: &cmd.Config[testConfig]{
				Name:        strings.Repeat("a", configNameMaxLen),
				Description: strings.Repeat("b", configDescriptionMaxLen+1),
			},
			ExpectedErr: errConfigDescriptionLen,
		},
		{
			Scenario: "execution nil",
			Adapter:  New[testConfig](),
			Config: &cmd.Config[testConfig]{
				Name:        strings.Repeat("a", configNameMaxLen),
				Description: strings.Repeat("b", configDescriptionMaxLen),
				Execute:     nil,
			},
			ExpectedErr: errConfigExecuteNil,
		},
	}

	for _, rowTest := range tableTest {
		rowTest := rowTest

		t.Run(rowTest.Scenario, func(t *testing.T) {
			t.Parallel()

			actualErr := rowTest.Adapter.Add(rowTest.Config)
			assert.Equal(t, fmt.Errorf("%s: %w", cmd.MODULE_NAME, rowTest.ExpectedErr).Error(), actualErr.Error())
		})
	}
}

func TestAdd_ConfigMaxLen(t *testing.T) {
	t.Parallel()

	adapter := New[testConfig]()

	for i := len(adapter.configs) + 1; i <= configMaxLen; i++ {
		actualErr := adapter.Add(&cmd.Config[testConfig]{
			Name:        "command_" + strconv.Itoa(i),
			Description: "description " + strconv.Itoa(i),
			Execute: func(app app.V1[testConfig]) error {
				return nil
			},
		})
		assert.Equal(t, nil, actualErr)
	}

	actualErr := adapter.Add(&cmd.Config[testConfig]{
		Name:        "command_" + strconv.Itoa(configMaxLen+1),
		Description: "description " + strconv.Itoa(configMaxLen+1),
		Execute: func(app app.V1[testConfig]) error {
			return nil
		},
	})
	assert.Equal(t, "cmd: configs must be equal or less than 10", actualErr.Error())
}

func TestAdd_AlreadyExist(t *testing.T) {
	t.Parallel()

	adapter := New[testConfig]()

	configA := &cmd.Config[testConfig]{
		Name:        "command A",
		Description: "description",
		Execute: func(app app.V1[testConfig]) error {
			return nil
		},
	}

	configB := &cmd.Config[testConfig]{
		Name:        configA.Name,
		Description: "description",
		Execute: func(app app.V1[testConfig]) error {
			return nil
		},
	}

	actualErr := adapter.Add(configA)
	assert.Equal(t, nil, actualErr)

	actualErr = adapter.Add(configB)
	assert.Equal(t, fmt.Errorf("%s: %w: %s", cmd.MODULE_NAME, errConfigAlreadyExist, configA.Name), actualErr)
}

func TestReplaceSpaceSplit(t *testing.T) {
	t.Parallel()

	actual := replaceSpaceSplit("")
	assert.Equal(t, []string{}, actual)

	actual = replaceSpaceSplit(" ")
	assert.Equal(t, []string{}, actual)

	actual = replaceSpaceSplit("a ")
	assert.Equal(t, []string{"a"}, actual)

	actual = replaceSpaceSplit(" a")
	assert.Equal(t, []string{"a"}, actual)

	actual = replaceSpaceSplit(" a b ")
	assert.Equal(t, []string{"a", "b"}, actual)

	actual = replaceSpaceSplit(" a  b ")
	assert.Equal(t, []string{"a", "b"}, actual)
}
