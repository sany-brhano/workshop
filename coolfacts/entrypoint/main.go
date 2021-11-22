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
            <img src="{{.Image}}" width="80%" />
       </article>
    {{ end }}
  <div>
  </body>
</html>`

func main() {
	factRepo := repo{
		facts: []fact{
			{Image: "pic1", Description: "DOGS PIC"},
			{Image: "pic2", Description: "CATS PIC"},
		},
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "PONG")
		if err != nil {
			errMessage := fmt.Sprintf("error writing response: %v", err)
			http.Error(w, errMessage, http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/facts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		if r.Method == http.MethodGet {
			tmpl, err := template.New("facts").Parse(newsTemplate)
			if err != nil {
				return
			}
			facts := factRepo.getAll()
			err = tmpl.Execute(w, facts)
			if err != nil {
				errMessage := fmt.Sprintf("error executing html: %v", err)
				http.Error(w, errMessage, http.StatusInternalServerError)
			}
		}

		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				errMessage := fmt.Sprintf("error reading response: %v", err)
				http.Error(w, errMessage, http.StatusInternalServerError)
			}
			var factPost fact
			err = json.Unmarshal(body, &factPost)
			if err != nil {
				fmt.Println("error Unmarshallling json slice:", err)
			}

			f := fact{
				Image:       factPost.Image,
				Description: factPost.Description,
			}
			factRepo.add(f)

			_, wErr := w.Write([]byte("SUCCESS"))
			if wErr != nil {
				fmt.Println("error writing SUCCESS: ", wErr)
			}
		}
	})

	log.Println("starting server")
	err := http.ListenAndServe(":9002", nil)
	if err != nil {
		fmt.Println("error listening server: ", err)
	}
}
