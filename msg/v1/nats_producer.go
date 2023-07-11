package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ampliway/way-lib-go/ctx"
	"github.com/ampliway/way-lib-go/helper/id"
	"github.com/ampliway/way-lib-go/helper/reflection"
	"github.com/ampliway/way-lib-go/msg"
	"github.com/iancoleman/strcase"
	"github.com/nats-io/nats.go"
)

var _ msg.ProducerV1 = (*NatsProducer)(nil)

type NatsProducer struct {
	nc *nats.Conn
	js nats.JetStreamContext
	id id.ID
}

func New(config *Config, id id.ID) (*NatsProducer, error) {
	if config == nil {
		return nil, fmt.Errorf("%s: %w", msg.MODULE_NAME, errConfigNull)
	}

	if config.NatsServers == "" {
		return nil, fmt.Errorf("%s: %w", msg.MODULE_NAME, errConfigServersEmpty)
	}

	nc, err := nats.Connect(config.NatsServers)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", msg.MODULE_NAME, errNATSConnect, err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", msg.MODULE_NAME, errJSConnect, err)
	}

	packageName := strings.ToLower(reflection.AppNamePkg())

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     packageName,
		Subjects: []string{fmt.Sprintf("%s.>", packageName)},
	})
	if err != nil {
		panic(err)
	}

	return &NatsProducer{
		nc: nc,
		js: js,
		id: id,
	}, nil
}

func (np *NatsProducer) Publish(ctx ctx.V1, m interface{}) error {
	subject := subjectName(m)
	expectedSubject := strings.ToLower(reflection.AppNamePkg()) + "."
	if !strings.HasPrefix(subject, expectedSubject) {
		return fmt.Errorf("%s: %w: receive: %s; expected: %s", msg.MODULE_NAME, errSubPrefix, subject, expectedSubject)
	}

	value, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("%s: %w: %s", msg.MODULE_NAME, errUnmarshal, subject)
	}

	header := nats.Header{}
	header.Add(msg.HEADER_X_TRACE_ID, ctx.TraceID())
	header.Add(msg.HEADER_X_MSG_ID, np.id.Random())

	_, err = np.js.PublishMsg(&nats.Msg{
		Subject: subject,
		Header:  header,
		Data:    value,
	})
	if err != nil {
		return fmt.Errorf("%s: %w: %s", msg.MODULE_NAME, errPublish, subject)
	}

	return err
}

func (np *NatsProducer) Shutdown() {
	np.nc.Drain()
}

func subjectName(msg interface{}) string {
	msgName := strcase.ToKebab(reflection.TypeName(msg))
	packageName := strings.ToLower(reflection.AppName(msg))
	return fmt.Sprintf("%s.%s", packageName, msgName)
}
