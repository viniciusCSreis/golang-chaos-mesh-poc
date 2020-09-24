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
	count = 5
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
	router.GET("/health", healthHandler)
	router.GET("/health/readiness", healthReadyHandler)

	log.Info().Msg("Hello from author manager")

	log.Info().Msg("Server: Running")
	log.Fatal().Err(http.ListenAndServe(":3000", router)).Msg("failed to start server!")

}

func healthHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func healthReadyHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	count++
	if count > 35 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
