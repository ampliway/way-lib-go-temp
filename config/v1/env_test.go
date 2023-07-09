package v1

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFieldNotSupportedConfig struct {
	FieldX uint64
}

type testConfig struct {
	Field1 string
	Field2 int
	Field3 bool
}

func TestNew(t *testing.T) {
	scenarios := []struct {
		Name        string
		Test        func() (any, error)
		Expected    interface{}
		ExpectedErr string
	}{
		{
			Name: "type_not_suported",
			Test: func() (any, error) {
				return New[string]()
			},
			Expected:    (*Env[string])(nil),
			ExpectedErr: "config: only structs are supported by config module",
		},
		{
			Name: "field_not_found",
			Test: func() (any, error) {
				return New[testConfig]()
			},
			Expected:    (*Env[testConfig])(nil),
			ExpectedErr: "config: environment variable not found: FIELD_1",
		},
		{
			Name: "field_not_supported",
			Test: func() (any, error) {
				os.Setenv("FIELD_X", "1")

				return New[testFieldNotSupportedConfig]()
			},
			Expected:    (*Env[testFieldNotSupportedConfig])(nil),
			ExpectedErr: "config: only \"string\", \"int\" and \"bool\" types are supported in struct",
		},
		{
			Name: "field_int_parse_failed",
			Test: func() (any, error) {
				os.Setenv("FIELD_1", "a")
				os.Setenv("FIELD_2", "b")
				os.Setenv("FIELD_3", "c")

				return New[testConfig]()
			},
			Expected:    (*Env[testConfig])(nil),
			ExpectedErr: "config: could not parse found value to integer: FIELD_2: value \"b\"",
		},
		{
			Name: "field_bool_parse_failed",
			Test: func() (any, error) {
				os.Setenv("FIELD_1", "a")
				os.Setenv("FIELD_2", "1")
				os.Setenv("FIELD_3", "c")

				return New[testConfig]()
			},
			Expected:    (*Env[testConfig])(nil),
			ExpectedErr: "config: could not parse found value to boolean: FIELD_3: value \"c\"",
		},
		{
			Name: "success_1",
			Test: func() (any, error) {
				os.Setenv("FIELD_1", "a")
				os.Setenv("FIELD_2", "1")
				os.Setenv("FIELD_3", "true")

				return New[testConfig]()
			},
			Expected:    &Env[testConfig]{value: &testConfig{Field1: "a", Field2: 1, Field3: true}},
			ExpectedErr: "",
		},
		{
			Name: "success_2",
			Test: func() (any, error) {
				os.Setenv("FIELD_1", "b")
				os.Setenv("FIELD_2", "-1")
				os.Setenv("FIELD_3", "false")

				return New[testConfig]()
			},
			Expected:    &Env[testConfig]{value: &testConfig{Field1: "b", Field2: -1, Field3: false}},
			ExpectedErr: "",
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario

		actualValue, actualErr := scenario.Test()

		assert.Equal(t, scenario.Expected, actualValue)

		if actualErr == nil {
			assert.Equal(t, scenario.ExpectedErr, "")
		} else {
			assert.Equal(t, scenario.ExpectedErr, actualErr.Error())
		}
	}
}

func TestGet(t *testing.T) {
	os.Setenv("FIELD_1", "b")
	os.Setenv("FIELD_2", "-1")
	os.Setenv("FIELD_3", "false")

	env, err := New[testConfig]()
	assert.Equal(t, nil, err)

	value := env.Get()

	assert.Equal(t, &testConfig{
		Field1: "b",
		Field2: -1,
		Field3: false,
	}, value)
}
