package books

import (
	"github.com/rs/zerolog/log"
	"plataform/pkg/provider/messaging"
)

type CreatorMessaging struct {
	m messaging.Subscriber
}

func (h CreatorMessaging) Handler() {
	go h.m.Subscribe(messaging.SubjectBuildBooK, func(m messaging.Message) {
		log.Info().Msgf("Msg %v", m)
	}, func(err error) {
		log.Error().Err(err)
	})
}

func NewCreatorMessaging(msg messaging.Subscriber) CreatorMessaging {
	return CreatorMessaging{msg}
}
