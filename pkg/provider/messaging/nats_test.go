package messaging

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"plataform/pkg/api"
	"plataform/pkg/provider"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func Test_natsMessenger_Subscribe(t *testing.T) {
	type in struct {
		nc    provider.NatsExecutor
		queue string
		s     Subject
		hm    MessageHandler
		he    ErrorHandler
	}

	teste1Err := errors.New("teste1Err not changed")
	wg1 := sync.WaitGroup{}
	wg1.Add(1)

	teste2Err := errors.New("teste2Err not changed")
	wg2 := sync.WaitGroup{}
	wg2.Add(1)

	teste3Err := errors.New("teste3Err not changed")
	wg3 := sync.WaitGroup{}
	wg3.Add(1)

	teste4Err := errors.New("teste4Err not changed")
	wg4 := sync.WaitGroup{}
	wg4.Add(1)

	teste5Err := errors.New("teste5Err not changed")
	wg5 := sync.WaitGroup{}
	wg5.Add(1)

	tests := []struct {
		name string
		in   in
		err  func() error
	}{
		{
			name: "Run with success",
			in: in{
				nc: NatsExecutorCustomMock{
					QueueSubscribeMock: func(subj, queue string, cb provider.NatsMsgHandler) (*nats.Subscription, error) {

						msg := Message{Content: []byte(`{"teste":"OK"}`)}
						msgEncode, _ := encode(msg)

						natsMsg := NatsMsgCustomMock{
							RespondMock: func(data []byte) error {
								return nil
							},
							DataMock: func() []byte {
								return msgEncode
							},
						}
						cb(&natsMsg)
						return nil, nil
					},
				},
				s:     "teste",
				queue: "server",
				he: func(err error) {
					if err != nil {
						teste1Err = errors.New("err should be nil")
					}
				},
				hm: func(m Message) {
					var result map[string]string
					if err := json.Unmarshal(m.Content, &result); err != nil {
						teste1Err = errors.New("fail to unmarshal result")
					}
					if c, exist := result["teste"]; !exist || c != "OK" {
						teste1Err = errors.New("result Message is not {'teste':'OK'}")
					} else {
						teste1Err = nil
					}
					wg1.Done()
				},
			},
			err: func() error {
				wg1.Wait()
				return teste1Err
			},
		},
		{
			name: "Call he when responde fail",
			in: in{
				nc: NatsExecutorCustomMock{
					QueueSubscribeMock: func(subj, queue string, cb provider.NatsMsgHandler) (*nats.Subscription, error) {

						msg := Message{Content: []byte(`{"teste":"OK"}`)}
						msgEncode, _ := encode(msg)

						natsMsg := NatsMsgCustomMock{
							RespondMock: func(data []byte) error {
								return errors.New("respond fail")
							},
							DataMock: func() []byte {
								return msgEncode
							},
						}
						cb(&natsMsg)
						return nil, nil
					},
				},
				s:     "teste",
				queue: "server",
				he: func(err error) {
					exp := errors.New("respond fail").Error()
					if err.Error() != exp {
						teste2Err = fmt.Errorf("error should be: %s\n", exp)
					} else {
						teste2Err = nil
					}
					wg2.Done()
				},
			},
			err: func() error {
				wg2.Wait()
				return teste2Err
			},
		},
		{
			name: "Call he when fail to decode base64",
			in: in{
				nc: NatsExecutorCustomMock{
					QueueSubscribeMock: func(subj, queue string, cb provider.NatsMsgHandler) (*nats.Subscription, error) {

						natsMsg := NatsMsgCustomMock{
							RespondMock: func(data []byte) error {
								return nil
							},
							DataMock: func() []byte {
								return []byte("this is no a base64")
							},
						}
						cb(&natsMsg)
						return nil, nil
					},
				},
				s:     "teste",
				queue: "server",
				he: func(err error) {
					exp := errors.New("illegal base64 data at input byte 4").Error()
					if err.Error() != exp {
						teste3Err = fmt.Errorf("error should be: %s\n", exp)
					} else {
						teste3Err = nil
					}
					wg3.Done()
				},
			},
			err: func() error {
				wg3.Wait()
				return teste3Err
			},
		},
		{
			name: "Call he when json decode fail",
			in: in{
				nc: NatsExecutorCustomMock{
					QueueSubscribeMock: func(subj, queue string, cb provider.NatsMsgHandler) (*nats.Subscription, error) {

						msgEncode := base64.StdEncoding.EncodeToString([]byte(`this is not a json`))

						natsMsg := NatsMsgCustomMock{
							RespondMock: func(data []byte) error {
								return nil
							},
							DataMock: func() []byte {
								return []byte(msgEncode)
							},
						}
						cb(&natsMsg)
						return nil, nil
					},
				},
				s:     "teste",
				queue: "server",
				he: func(err error) {
					exp := errors.New("invalid character 'h' in literal true (expecting 'r')").Error()
					if err.Error() != exp {
						teste4Err = fmt.Errorf("error should be: %s\n", exp)
					} else {
						teste4Err = nil
					}
					wg4.Done()
				},
			},
			err: func() error {
				wg4.Wait()
				return teste4Err
			},
		},
		{
			name: "Call he when QueueSubscribeMock return err",
			in: in{
				nc: NatsExecutorCustomMock{
					QueueSubscribeMock: func(subj, queue string, cb provider.NatsMsgHandler) (*nats.Subscription, error) {

						return nil, errors.New("error on QueueSubscribeMock")
					},
				},
				s:     "teste",
				queue: "server",
				he: func(err error) {
					if err.Error() != errors.New("error on QueueSubscribeMock").Error() {
						teste5Err = errors.New("error should be 'error on QueueSubscribeMock'")
					} else {
						teste5Err = nil
					}
					wg5.Done()
				},
			},
			err: func() error {
				wg5.Wait()
				return teste5Err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewMessenger(tt.in.nc, tt.in.queue)
			n.Subscribe(tt.in.s, tt.in.hm, tt.in.he)
			if err := tt.err(); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_natsMessenger_PublishSync(t *testing.T) {
	type in struct {
		org   api.Organization
		nc    provider.NatsExecutor
		queue string
		s     Subject
		m     Message
	}

	tests := []struct {
		name string
		in   in
		err  error
	}{
		{
			name: "Publish with success",
			in: in{
				org: "zup",
				nc: NatsExecutorCustomMock{
					RequestMock: func(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
						msg, err := decode(data)
						if err != nil {
							return nil, err
						}

						if h, exist := msg.Header["teste"]; h != "this_is_a_header" || !exist {
							return nil, errors.New("header should have 'teste: this_is_a_header'")
						}
						if h, exist := msg.Header[OrganizationHeader]; h != "zup" || !exist {
							return nil, errors.New("header should have 'OrganizationHeader: zup'")
						}
						content, err := base64.StdEncoding.DecodeString(string(msg.Content))
						if err != nil {
							return nil, err
						}
						if string(content) != `{"teste2":"this_is_fine"}` {
							return nil, errors.New(`content should be : {"teste2":"this_is_fine"}`)
						}
						return nil, nil
					},
				},
				s:     "teste",
				queue: "server",
				m: Message{
					Header:  map[string]string{"teste": "this_is_a_header"},
					Content: []byte(base64.StdEncoding.EncodeToString([]byte(`{"teste2":"this_is_fine"}`))),
				},
			},
			err: nil,
		},
		{
			name: "Publish without header",
			in: in{
				org: "zup",
				nc: NatsExecutorCustomMock{
					RequestMock: func(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
						msg, err := decode(data)
						if err != nil {
							return nil, err
						}

						if h, exist := msg.Header[OrganizationHeader]; h != "zup" || !exist {
							return nil, errors.New("header should have 'OrganizationHeader: zup'")
						}
						content, err := base64.StdEncoding.DecodeString(string(msg.Content))
						if err != nil {
							return nil, err
						}
						if string(content) != `{"teste2":"this_is_fine"}` {
							return nil, errors.New(`content should be : {"teste2":"this_is_fine"}`)
						}
						return nil, nil
					},
				},
				s:     "teste",
				queue: "server",
				m: Message{
					Content: []byte(base64.StdEncoding.EncodeToString([]byte(`{"teste2":"this_is_fine"}`))),
				},
			},
			err: nil,
		},
		{
			name: "Return err when timeout",
			in: in{
				org: "zup",
				nc: NatsExecutorCustomMock{
					RequestMock: func(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
						return nil, errors.New("timeout")
					},
					LastErrorMock: func() error {
						return errors.New("last error is timeout")
					},
				},
				s:     "teste",
				queue: "server",
				m: Message{
					Header:  map[string]string{"teste": "this_is_a_header"},
					Content: []byte(base64.StdEncoding.EncodeToString([]byte(`{"teste2":"this_is_fine"}`))),
				},
			},
			err: errors.New("timeout"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewMessenger(tt.in.nc, tt.in.queue)
			err := n.PublishSync(tt.in.org, tt.in.s, tt.in.m)

			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("%s:\n expected: %v\n recevided: %v", tt.name, tt.err, err)
				}
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("%s:\n expected: %v\n recevided: %v", tt.name, tt.err, err)
			}
		})
	}
}
