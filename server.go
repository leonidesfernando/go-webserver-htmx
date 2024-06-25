package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
)


var 	(
	films = map[string][]Film{
		"Films":{
			{Title: "The Godfather", Director: "Francis Ford Coppola"},
			{Title: "Pulp Fiction", Director: "Quentin Tarantino"},
			{Title: "Jurassic Park", Director: "Steven Spielberg"},
			{Title: " Back to the future", Director: " Robert Zemeckis"},
		},
	}

	mu    sync.Mutex
)

type Result struct{
	Message string
}

func main() {
	const port = 8088;

	fmt.Println("We are ready to serve you!")
	fmt.Printf("Running on http://localhost:%d/\n", port)

	endPoints(films)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func endPoints(films map[string][]Film){
	home(films)
	addFilm(films)
}


func addFilm(films map[string][]Film){
	addFilmHandler := func (writer http.ResponseWriter, request *http.Request)  {
		title := request.PostFormValue("title")
		director := request.PostFormValue("director")
		newItem := Film{Title: title,Director: director}

		list := films["Films"]

		mu.Lock()
    defer mu.Unlock()

		tmpl := template.Must(template.ParseFiles("index.html"))
		if exists(list, newItem) {

			result := Result{Message: "This movie and director were added already."}
			tmpl.ExecuteTemplate(writer, "message", result)
      
		}else {

			list = append(list, newItem)
			films["Films"] = list
			tmpl.ExecuteTemplate(writer, "film-list-element",newItem)
		}
	}
	http.HandleFunc("/add-film/", addFilmHandler)
}

func home(films map[string][]Film){
	homeHandlder := func (writer http.ResponseWriter, request *http.Request)  {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(writer, films)
	}		
	http.HandleFunc("/", homeHandlder)
}

// Function to check if a film exists in the slice
func exists(films []Film, target Film) bool {
	for _, film := range films {
		if film == target {
			return true
		}
	}
	return false
}
