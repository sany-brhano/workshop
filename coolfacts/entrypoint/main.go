package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var newsTemplate = `<!DOCTYPE html>
<html>
  <head><style>/* copy coolfacts/styles.css for some color 🎨*/</style></head>
  <body>
  <h1>Facts List</h1>
  <div>
    {{ range . }}
       <article>
            <h3>{{.Description}}</h3>
            <img src="{{.Image}}" width="100%" />
       </article>
    {{ end }}
  <div>
  </body>
</html>`

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

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "PONG")
		if err != nil {
			errMessage := fmt.Sprintf("error writing response: %v", err)
			http.Error(w, errMessage, http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/facts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.New("facts").Parse(newsTemplate)
			if err != nil {
				fmt.Println("error parsing facts to html: ", err)
			}
			facts := factsRepo.getAll()
			err = tmpl.Execute(w, facts)
			if err != nil {
				errMessage := fmt.Sprintf("error executing html: %v", err)
				http.Error(w, errMessage, http.StatusInternalServerError)
			}
		}
		if r.Method == http.MethodPost {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				errMessage := fmt.Sprintf("error reading response: %v", err)
				http.Error(w, errMessage, http.StatusInternalServerError)
			}
			var factPost fact
			err = json.Unmarshal(b, &factPost)
			if err != nil {
				fmt.Println("error unmarshalling json: ", err)
			}

			f := fact{
				Image:       factPost.Image,
				Description: factPost.Description,
			}
			factsRepo.add(f)
			_, wErr := w.Write([]byte("SUCCESS"))
			if wErr != nil {
				errMessage := fmt.Sprintf("error writing response: %v", err)
				http.Error(w, errMessage, http.StatusInternalServerError)
			}
		}
	})
	log.Println("starting server")
	err := http.ListenAndServe(":9002", nil)
	if err != nil {
		fmt.Println("error listening server: ", err)
	}
}
