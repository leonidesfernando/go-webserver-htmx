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
			{Title: "Steven Spielberg", Director: "Jurassic Park"},
			{Title: " Back to the future", Director: " Robert Zemeckis"},
		},
	}

	mu    sync.Mutex
)

func main() {
	const port = 8088;

	/*
	file,err := os.Open("data.json")	

	if(err != nil){
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if(err != nil){
		log.Fatal(err)
	}

	var data Films
	json.Unmarshal(bytes, &data)

	fmt.Println(data)

	for _, item := range data.films {
		fmt.Println(item.Title + " - " + item.Director	)
	} */

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

		if exists(list, newItem) {
			//fmt.Println("Já existe: ", newItem)
			writer.WriteHeader(http.StatusBadRequest)
      fmt.Fprint(writer, "This movie and director were added already.")
      return
		}
		list = append(list, newItem)
		fmt.Println(list)
		films["Films"] = list
		/*if (value != nil) && strings.ToLower(value) == strings.ToLower(director){

		}*/


		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(writer, "film-list-element",newItem)
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
	fmt.Println(target)
	for _, film := range films {
		fmt.Println(film)
		if film == target {
			fmt.Println("Achou")
			return true
		}
	}
	fmt.Println("Não achou")
	return false
}