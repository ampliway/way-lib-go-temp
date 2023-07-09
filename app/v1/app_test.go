package v1

import (
	"os"
	"testing"

	configV1 "github.com/ampliway/way-lib-go/config/v1"
	"github.com/stretchr/testify/assert"
)

type appConfig struct {
	FieldA string
}

func TestNew(t *testing.T) {
	os.Setenv("FIELD_A", "value a")

	env, err := configV1.New[appConfig]()
	assert.Nil(t, err)

	app, err := New[appConfig]()
	assert.Nil(t, err)

	assert.Equal(t, app, &App[appConfig]{config: env})

	config := app.Config().Get()
	assert.Equal(t, &appConfig{
		FieldA: "value a",
	}, config)
}
