package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMock(t *testing.T) {
	scenarios := []struct {
		Name        string
		Test        func() (any, error)
		Expected    interface{}
		ExpectedErr string
	}{
		{
			Name: "success",
			Test: func() (any, error) {
				return NewMock[testConfig](testConfig{Field1: "a", Field2: 2, Field3: true})
			},
			Expected:    &Mock[testConfig]{value: &testConfig{Field1: "a", Field2: 2, Field3: true}},
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

func TestMockGet(t *testing.T) {
	env, err := NewMock[testConfig](testConfig{Field1: "a", Field2: 2, Field3: true})
	assert.Equal(t, nil, err)

	value := env.Get()

	assert.Equal(t, &testConfig{
		Field1: "a",
		Field2: 2,
		Field3: true,
	}, value)
}
