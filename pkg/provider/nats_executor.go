package provider

import (
	"fmt"
	"io/ioutil"
	"plataform/pkg/api"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type natsExecutor struct {
	nc *nats.Conn
}

func NewNatsExecutor(config api.Configs) (NatsExecutor, error) {

	natsUrl := fmt.Sprintf("nats://%s:4222", config.AppConfigs.NatsConfig.Host)

	if config.AppConfigs.NatsConfig.User != "" {
		token, err := ioutil.ReadFile("/var/run/secrets/nats.io/token")
		if err != nil {
			return natsExecutor{}, err
		}
		natsUrl = fmt.Sprintf(
			"nats://%s:%s@%s:4222",
			config.AppConfigs.NatsConfig.User,
			string(token),
			config.AppConfigs.NatsConfig.Host,
		)
	}

	nc, err := nats.Connect(
		natsUrl,
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Err(err).Msg("nats client disconnected")
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Info().Msg("nats client reconnected")
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Info().Msg("nats client closed")
		}),
	)
	if err != nil {
		log.Err(err).Msg("Fail to connect to nats.")
	}
	return natsExecutor{
		nc: nc,
	}, err
}

func (n natsExecutor) QueueSubscribe(subj, queue string, cb NatsMsgHandler) (*nats.Subscription, error) {

	return n.nc.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		proxy := natsMsg{
			msg: msg,
		}
		cb(&proxy)
	})
}

func (n natsExecutor) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	return n.nc.Request(subj, data, timeout)
}

func (n natsExecutor) LastError() error {
	return n.nc.LastError()
}

type natsMsg struct {
	msg *nats.Msg
}

func (m natsMsg) Respond(data []byte) error {
	return m.msg.Respond(data)
}

func (m natsMsg) Data() []byte {
	return m.msg.Data
}
