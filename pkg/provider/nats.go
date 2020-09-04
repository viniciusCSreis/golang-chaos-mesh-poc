package provider

import (
	"time"

	"github.com/nats-io/nats.go"
)

type NatsMsg interface {
	Respond(data []byte) error
	Data() []byte
}

type NatsMsgHandler func(msg NatsMsg)

type NatsExecutor interface {
	QueueSubscribe(subj, queue string, cb NatsMsgHandler) (*nats.Subscription, error)
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
	LastError() error
}
