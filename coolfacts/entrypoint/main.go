package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	factsRepo := repo{
		facts: []fact{
			{Image: "pic1", Description: "DOGS PIC"},
		},
	}
	factsRepo.add(fact{
		Image:       "https://images2.minutemediacdn.com/image/upload/v1556645500/shape/cover/entertainment/D5aliXvWsAEcYoK-fe997566220c082b98030508e654948e.jpg",
		Description: "Did you know sonic is a hedgehog?!",
	})

	http.HandleFunc("/facts", func(w http.ResponseWriter, r *http.Request) {
		allFacts := factsRepo.getAll()
		bs, err := json.Marshal(allFacts)
		if err != nil {
			errMessage := fmt.Sprintf(" marshal error")
			http.Error(w, errMessage, http.StatusInternalServerError)
		}
		_, err = w.Write(bs)
		if err != nil {
			errMessage := fmt.Sprintf("error writing response: %v", err)
			http.Error(w, errMessage, http.StatusInternalServerError)
		}

	})

	log.Println("starting server")
	err := http.ListenAndServe(":9002", nil)
	if err != nil {
		fmt.Println("error listening server: ", err)
	}
}
