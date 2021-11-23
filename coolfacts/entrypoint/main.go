package main

import (
	"context"
	"fmt"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/fact"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/mentalfloss"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/myhttp"
	"log"
	"net/http"
	"time"
)

var factRepo fact.Repo

func mfUpdate() error {
	mfFacts, err := mentalfloss.Mentalfloss{}.Facts()
	if err != nil {
		log.Fatal(err)
	}
	factRepo = fact.Repo{}
	for _, val := range mfFacts {
		factRepo.Add(val)
	}
	return nil
}

func updateFactsWithTicker(ctx context.Context, updateFunc func() error) {
	ticker := time.NewTicker(5 * time.Millisecond)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := updateFunc()
				if err != nil {
					fmt.Println("error executing html: ", err)
				}
				fmt.Println("updating mentalfloss facts")
			}
		}
	}(ctx)
}

func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	updateFactsWithTicker(ctx, mfUpdate)

	handlerer := myhttp.FactsHandler{
		FactRepo: &factRepo,
	}
	http.HandleFunc("/ping", handlerer.Ping)
	http.HandleFunc("/facts", handlerer.Facts)

	log.Println("starting server")
	err := http.ListenAndServe(":9002", nil)
	if err != nil {
		fmt.Println("error listening server: ", err)
	}
}
