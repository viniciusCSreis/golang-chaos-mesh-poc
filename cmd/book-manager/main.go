package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"plataform/pkg/api"
	"plataform/pkg/books"
	"plataform/pkg/git"
	"plataform/pkg/git/github"
	"plataform/pkg/provider"
	"plataform/pkg/provider/messaging"
	"time"
)

var (
	appName = api.AppName("author-manager")
)

func main() {

	config := api.Configs{
		AppConfigs: api.AppConfigs{
			NatsConfig: api.NatsConfig{
				Host: "example-nats-cluster",
				User: "",
			},
		},
	}

	//Messaging
	nats, err := provider.NewNatsExecutor(config)
	if err != nil {
		os.Exit(1)
	}
	natsMessaging := messaging.NewMessenger(nats, appName.String())

	//Github
	githubRepo := github.NewRepoManager(&http.Client{Timeout: time.Second * 2})

	repoProviders := git.RepoProviders{}
	repoProviders.Add(git.GitHub, git.Git{Repos: githubRepo, NewRepoInfo: github.NewRepoInfo})

	bookListerHTTP := books.NewListerHTTP()
	bookCreatorMessaging := books.NewCreatorMessaging(natsMessaging, repoProviders)

	bookCreatorMessaging.Handler()

	router := httprouter.New()
	router.GET("/books", bookListerHTTP.Handler())
	router.GET("/health", healthHandler)
	router.GET("/health/readiness", healthHandler)

	log.Info().Msg("Server: Running")
	log.Fatal().Err(http.ListenAndServe(":3000", router)).Msg("failed to start server!")

}

func healthHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	_ , err := http.Get("http://author-manager:3000/authors")
	if err != nil {
		log.Error().Err(err).Msg("fail to get authors")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
