package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

	r, err := api.NewRouter(os.Getenv("SKEY"), os.Getenv("POSTGRES_DNS"))
	if err != nil {
		log.Error().Msgf("Bootstrap router: %s", err.Error())
		return
	}

	http.ListenAndServe(":"+os.Getenv("PORT"), r)

}
