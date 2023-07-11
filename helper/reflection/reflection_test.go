package reflection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeName(t *testing.T) {
	t.Parallel()

	stringPointer := "value1"

	rows := []struct {
		Scenario string
		Input    interface{}
		Expected string
	}{
		{"null", nil, ""},
		{"string", "value1", "string"},
		{"string_pointer", &stringPointer, "string"},
		{"struct", forTest{}, "forTest"},
		{"struct_pointer", &forTest{}, "forTest"},
	}

	for _, testData := range rows {
		testData := testData

		t.Run(testData.Scenario, func(t *testing.T) {
			t.Parallel()

			result := TypeName(testData.Input)
			assert.Equal(t, testData.Expected, result)
		})
	}
}

func TestAppName(t *testing.T) {
	t.Parallel()

	rows := []struct {
		Scenario string
		Input    interface{}
		Expected string
	}{
		{"success", &forTest{}, "way-lib-go"},
	}

	for _, testData := range rows {
		testData := testData

		t.Run(testData.Scenario, func(t *testing.T) {
			t.Parallel()

			result := AppName(testData.Input)
			assert.Equal(t, testData.Expected, result)
		})
	}
}

func TestAppNamePkg(t *testing.T) {
	t.Parallel()

	rows := []struct {
		Scenario string
		Expected string
	}{
		{"success_1", "way-lib-go"},
	}

	for _, testData := range rows {
		testData := testData

		t.Run(testData.Scenario, func(t *testing.T) {
			t.Parallel()

			result := AppNamePkg()
			assert.Equal(t, testData.Expected, result)
		})
	}
}

type forTest struct {
	Field1 string
	Field2 string
}
