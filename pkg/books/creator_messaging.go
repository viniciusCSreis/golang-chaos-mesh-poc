package books

import (
	"github.com/rs/zerolog/log"
	"plataform/pkg/git"
	"plataform/pkg/provider/messaging"
)

type CreatorMessaging struct {
	m     messaging.Subscriber
	repos git.RepoProviders
}

func (h CreatorMessaging) Handler() {
	go h.m.Subscribe(messaging.SubjectBuildBooK, func(m messaging.Message) {
		book, err := h.repos[git.GitHub].Repos.Zipball(
			h.repos[git.GitHub].NewRepoInfo("http://github.com/zupIt/ritchie-formulas", ""),
			"2.5.0",
		)
		if err != nil {
			log.Error().Err(err).Msg("Fail to download book")
		}

		log.Info().Msgf("Book:%v\n", book)
		log.Info().Msgf("Msg %v", m)
	}, func(err error) {
		log.Error().Err(err)
	})
}

func NewCreatorMessaging(msg messaging.Subscriber, repos git.RepoProviders) CreatorMessaging {
	return CreatorMessaging{m: msg, repos: repos}
}
