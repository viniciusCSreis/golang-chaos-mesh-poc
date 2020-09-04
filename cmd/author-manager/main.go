package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"plataform/pkg/api"
	"plataform/pkg/authors"
	"plataform/pkg/provider"
	"plataform/pkg/provider/messaging"
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

	authorListerHTTP := authors.NewListerHTTP()
	authorCreatorHTTP := authors.NewCreatorHTTP(natsMessaging)

	router := httprouter.New()
	router.GET("/authors", authorListerHTTP.Handler())
	router.POST("/authors", authorCreatorHTTP.Handler())

	log.Info().Msg("Hello from author manager")

	log.Info().Msg("Server: Running")
	log.Fatal().Err(http.ListenAndServe(":3000", router)).Msg("failed to start server!")

}
