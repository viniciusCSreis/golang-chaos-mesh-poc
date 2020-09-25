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
		lastTag, err := h.repos[git.GitHub].Repos.LatestTag(
			h.repos[git.GitHub].NewRepoInfo("http://github.com/zupIt/ritchie-formulas", ""),
		)
		if err != nil {
			log.Error().Err(err).Msg("Fail to get last tag")
		}

		log.Info().Msgf("Book:%v\n", lastTag)
		log.Info().Msgf("Msg %v", m)
	}, func(err error) {
		log.Error().Err(err)
	})
}

func NewCreatorMessaging(msg messaging.Subscriber, repos git.RepoProviders) CreatorMessaging {
	return CreatorMessaging{m: msg, repos: repos}
}
