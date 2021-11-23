package main

import (
	"fmt"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/fact"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/myhttp"
	"log"
	"net/http"
)

func main() {
	factsRepo := fact.Repo{}
	handlerer := myhttp.FactsHandler{
		FactRepo: &factsRepo,
	}
	http.HandleFunc("/ping", handlerer.Ping)

	log.Println("starting server")
	err := http.ListenAndServe(":9002", nil)
	if err != nil {
		fmt.Println("error listening server: ", err)
	}
}
