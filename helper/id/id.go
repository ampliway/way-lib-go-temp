package id

import (
	"strings"

	"github.com/oklog/ulid/v2"
)

var _ ID = (*Adapter)(nil)

type Adapter struct{}

func New() *Adapter {
	return new(Adapter)
}

func (a *Adapter) Random() string {
	return strings.ToLower(ulid.Make().String())
}
