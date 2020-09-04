package messaging

import (
	"encoding/base64"
	"encoding/json"
	"plataform/pkg/api"
	"plataform/pkg/provider"
	"time"

	"github.com/rs/zerolog/log"
)

var ackMsg = []byte("OK")

type natsMessenger struct {
	nc    provider.NatsExecutor
	queue string
}

func NewMessenger(nc provider.NatsExecutor, queue string) Messenger {
	return natsMessenger{nc: nc, queue: queue}
}

func (n natsMessenger) PublishSync(org api.Organization, s Subject, m Message) error {
	c := make(chan error)
	n.Publish(org, s, m, func(err error) {
		c <- err
	})
	return <-c
}

func (n natsMessenger) Subscribe(s Subject, hm MessageHandler, he ErrorHandler) {
	if _, err := n.nc.QueueSubscribe(string(s), n.queue, func(m provider.NatsMsg) {

		err := m.Respond(ackMsg)
		if err != nil {
			log.Err(err).Msg("Fail to Respond")
			he(err)
			return
		}
		go func() {
			msg, err := decode(m.Data())
			if err != nil {
				he(err)
				return
			}

			hm(msg)
			he(nil)
		}()

	}); err != nil {
		log.Err(err).Msgf("Fail to subscribe to subject: %s", string(s))
		he(err)
	}
}

func (n natsMessenger) Publish(org api.Organization, s Subject, m Message, he ErrorHandler) {

	go func() {
		if m.Header == nil {
			m.Header = map[string]string{}
		}
		m.Header[OrganizationHeader] = string(org)

		mBase64, err := encode(m)
		if err != nil {
			he(err)
			return
		}
		_, err = n.nc.Request(string(s), mBase64, 5*time.Second)
		if err != nil {
			if err := n.nc.LastError(); err != nil {
				log.Err(err).Msg("Last Request err")
			}
			log.Err(err).Msg("Fail Request make a request.")
			he(err)

			return
		}
		he(nil)
	}()
}

func encode(m Message) ([]byte, error) {
	mJson, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	mBase64 := base64.StdEncoding.EncodeToString(mJson)
	return []byte(mBase64), nil
}

func decode(data []byte) (Message, error) {
	mBase64, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		log.Err(err).Msg("Fail to decode base64 message")
		return Message{}, err
	}
	var mJson Message
	err = json.Unmarshal(mBase64, &mJson)
	if err != nil {
		log.Err(err).Msg("Fail to decode json message")
		return Message{}, err
	}
	return mJson, nil
}
