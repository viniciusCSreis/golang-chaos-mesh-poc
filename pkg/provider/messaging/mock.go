package messaging

import (
	"plataform/pkg/provider"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsExecutorCustomMock struct {
	QueueSubscribeMock func(subj, queue string, cb provider.NatsMsgHandler) (*nats.Subscription, error)
	RequestMock        func(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
	LastErrorMock      func() error
}

func (n NatsExecutorCustomMock) QueueSubscribe(subj, queue string, cb provider.NatsMsgHandler) (*nats.Subscription, error) {
	return n.QueueSubscribeMock(subj, queue, cb)
}

func (n NatsExecutorCustomMock) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	return n.RequestMock(subj, data, timeout)
}

func (n NatsExecutorCustomMock) LastError() error {
	return n.LastErrorMock()
}

type NatsMsgCustomMock struct {
	RespondMock func(data []byte) error
	DataMock    func() []byte
}

func (n NatsMsgCustomMock) Respond(data []byte) error {
	return n.RespondMock(data)
}

func (n NatsMsgCustomMock) Data() []byte {
	return n.DataMock()
}
