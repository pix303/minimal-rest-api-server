package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pix303/minimal-rest-api-server/pkg/api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error on retrive enviroment vars")
	}

	log.Println("hello minimal server is rolling!")

	r, err := api.NewRouter(os.Getenv("SKEY"), os.Getenv("POSTGRES_DNS"))
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":"+os.Getenv("PORT"), r)

}
