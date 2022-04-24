package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/pix303/minimal-rest-api-server/pkg/api"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Msgf("Loading env: %s", err.Error())
		return
	}

	log.Info().Msg("Hello minimal server api!")

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), os.Getenv("GITHUB_CALLBACK")),
	)

	r, err := api.NewRouter(os.Getenv("SKEY"), os.Getenv("GOTH_SESSION_SECRET"), os.Getenv("SESSION_SECRET"), os.Getenv("POSTGRES_DNS"))
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	if err != nil {
		log.Error().Msgf("Bootstrap router: %s", err.Error())
		return
	}

	http.ListenAndServe(":"+os.Getenv("PORT"), loggedRouter)
}
