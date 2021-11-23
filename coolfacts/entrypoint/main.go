package main

import (
	"context"
	"fmt"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/fact"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/inmem"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/myhttp"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/providers/mentalfloss"
	"log"
	"net/http"
)

func main() {
	factsRepo := inmem.NewFactRepository()
	mentalflossProvider := mentalfloss.NewProvider()
	service := fact.NewService(factsRepo, mentalflossProvider)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	service.UpdateFactsWithTicker(ctx, service.UpdateFacts)

	handlerer := myhttp.NewFactsHandler(factsRepo)
	http.HandleFunc("/ping", handlerer.Ping)
	http.HandleFunc("/facts", handlerer.Facts)

	log.Println("starting server")
	err := http.ListenAndServe(":9002", nil)
	if err != nil {
		fmt.Println("error listening server: ", err)
	}
}
