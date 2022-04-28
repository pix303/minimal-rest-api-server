package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/pix303/minimal-rest-api-server/pkg/api"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/authboss/v3"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Msgf("Loading env: %s", err.Error())
		return
	}

	log.Info().Msg("Hello minimal server api!")

	ab, err := api.SetupAuthConfig()
	if err != nil {
		log.Error().Err(err).Msg("fail to setup authorization")
	}

	ab.Config.Modules.OAuth2Providers = map[string]authboss.OAuth2Provider{
		"github": {
			OAuth2Config: &oauth2.Config{
				ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
				ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
				Scopes:       []string{"profile", "email"},
				Endpoint:     github.Endpoint,
			},
		},
	}

	if err := ab.Init(); err != nil {
		log.Error().Err(err).Msg("fail to init authorization")
	}

	r, err := api.NewRouter(os.Getenv("SKEY"), os.Getenv("POSTGRES_DNS"), ab)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	if err != nil {
		log.Error().Msgf("Bootstrap router: %s", err.Error())
		return
	}

	http.ListenAndServe(":"+os.Getenv("PORT"), loggedRouter)
}
