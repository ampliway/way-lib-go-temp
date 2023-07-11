package v1

import (
	"github.com/ampliway/way-lib-go/ctx"
	"github.com/ampliway/way-lib-go/helper/id"
)

var (
	_ ctx.V1 = (*Ctx)(nil)
)

type Ctx struct {
	traceID string
}

func New(id id.ID) *Ctx {
	return &Ctx{
		traceID: id.Random(),
	}
}

func (c *Ctx) TraceID() string {
	return c.traceID
}
