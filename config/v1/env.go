package v1

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"sync"

	"github.com/ampliway/way-lib-go/config"
	"github.com/iancoleman/strcase"
)

var (
	_                      config.V1[any] = (*Env[any])(nil)
	errGenericNotSupported                = errors.New("only structs are supported by config module")
	errFieldNotSupported                  = errors.New("only \"string\", \"int\" and \"bool\" types are supported in struct")
	errEnvNotFound                        = errors.New("environment variable not found")
	errEnvIntParse                        = errors.New("could not parse found value to integer")
	errEnvBoolParse                       = errors.New("could not parse found value to boolean")
)

type Env[T any] struct {
	value *T
}

var (
	once sync.Once
	args map[string]string
)

func New[T any]() (*Env[T], error) {
	once.Do(func() {
		args = map[string]string{}

		flag.VisitAll(func(f *flag.Flag) {
			args[f.Name] = f.Value.String()
		})
	})

	value := *new(T)
	if reflect.ValueOf(value).Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s: %w", config.MODULE_NAME, errGenericNotSupported)
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(&value).Elem()

	overideValues := map[string]string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		envName := strcase.ToScreamingSnake(f.Name)

		if value, exist := args[envName]; exist {
			overideValues[envName] = value
		}
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		envName := strcase.ToScreamingSnake(f.Name)

		envValue, envExist := overideValues[envName]
		if !envExist {
			envValue, envExist = os.LookupEnv(envName)
		}

		if !envExist {
			return nil, fmt.Errorf("%s: %w: %s", config.MODULE_NAME, errEnvNotFound, envName)
		}

		field := v.FieldByName(f.Name)

		switch f.Type.String() {
		case "string":
			field.SetString(envValue)
		case "int":
			envValueInt, err := strconv.Atoi(envValue)
			if err != nil {
				return nil, fmt.Errorf("%s: %w: %s: value \"%s\"", config.MODULE_NAME, errEnvIntParse, envName, envValue)
			}
			field.SetInt(int64(envValueInt))
		case "bool":
			envValueBool, err := strconv.ParseBool(envValue)
			if err != nil {
				return nil, fmt.Errorf("%s: %w: %s: value \"%s\"", config.MODULE_NAME, errEnvBoolParse, envName, envValue)
			}
			field.SetBool(envValueBool)
		default:
			return nil, fmt.Errorf("%s: %w", config.MODULE_NAME, errFieldNotSupported)
		}
	}

	return &Env[T]{
		value: &value,
	}, nil
}

func (e *Env[T]) Get() *T {
	return e.value
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
