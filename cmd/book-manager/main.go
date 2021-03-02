package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"plataform/pkg/api"
	"plataform/pkg/books"
	"plataform/pkg/provider"
	"plataform/pkg/provider/messaging"
)

var (
	appName = api.AppName("book-manager")
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

	bookListerHTTP := books.NewListerHTTP()
	bookCreatorMessaging := books.NewCreatorMessaging(natsMessaging)

	bookCreatorMessaging.Handler()

	router := httprouter.New()
	router.GET("/books", bookListerHTTP.Handler())
	router.GET("/health", healthHandler)

	log.Info().Msg("Server: Running")
	log.Fatal().Err(http.ListenAndServe(":3000", router)).Msg("failed to start server!")

}

func healthHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
