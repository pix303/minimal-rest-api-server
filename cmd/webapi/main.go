package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"

	"github.com/pix303/minimal-rest-api-server/pkg/api"
	"github.com/rs/zerolog/log"

	"github.com/pix303/minimal-rest-api-server/pkg/auth"
)

func main() {

	if os.Getenv("POSTGRES_DNS") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Error().Msgf("Loading env: %s", err.Error())
			return
		}
	}

	log.Info().Msg("Hello minimal api!")

	providers := make(map[string]auth.ProviderKeys)
	providers["github"] = auth.ProviderKeys{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Callback:     os.Getenv("GITHUB_CALLBACK"),
	}

	auth.InitOauth(providers, os.Getenv("OAUTH_SESSION_SECRET"))

	r, err := api.NewRouter(os.Getenv("POSTGRES_DNS"))

	if err != nil {
		log.Error().Msgf("Bootstrap router: %s", err.Error())
		return
	}

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	http.ListenAndServe(":"+os.Getenv("PORT"), loggedRouter)
}
