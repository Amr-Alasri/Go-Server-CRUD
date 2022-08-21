package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

func (m Movie) isEmpty() bool {
	return m.Isbn == ""
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

// function rourt

// GetMovies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(movies)
}

// GetOneMovie

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")

	// get the parameter from the request
	parms := mux.Vars(r)

	// loop over all movies , if found the movie return it otherwise retrun message('not found')
	for _, value := range movies {
		if value.ID == parms["id"] {
			json.NewEncoder(w).Encode(value)
			return
		}
	}

	// if no movie with that id this message well return
	json.NewEncoder(w).Encode("there is no movie with this id")
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {

	// set the content-type
	w.Header().Set("Content-Type", "Application/json")

	// get the parameter from the request
	parms := mux.Vars(r)

	// loop over all movies , if found the movie delete it otherwise retrun message('not found')
	for index, value := range movies {
		if value.ID == parms["id"] {
			// delete the movies
			movies = append(movies[:index], movies[index+1:]...)
			json.NewEncoder(w).Encode(movies)

			return
		}
	}

	// if no movie with that id this message well return
	json.NewEncoder(w).Encode("there is no movie with this id")

}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		json.NewEncoder(w).Encode("err:Body is nil")
		return
	}

	var movie Movie
	json.NewDecoder(r.Body).Decode(&movie)

	if movie.isEmpty() {
		json.NewEncoder(w).Encode("err:Body is Empty")
		return
	}
	rand.Seed(time.Now().UnixMilli())
	rNumber := rand.Intn(1000)
	movie.ID = strconv.Itoa(rNumber)

	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// if Body is nil
	if r.Body == nil {
		json.NewEncoder(w).Encode("err:Body is nil")
		return
	}

	// if Body is empty like {}

	parms := mux.Vars(r)

	//parse body to movie type

	var movie Movie

	json.NewDecoder(r.Body).Decode(&movie)

	for i, v := range movies {
		if v.ID == parms["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}

	json.NewEncoder(w).Encode("no Movie with this id:" + parms["id"])

}

func main() {
	r := mux.NewRouter()

	movies = append(movies,
		Movie{ID: "1", Isbn: "42351", Title: "Movie1", Director: &Director{FirstName: "Amr", LastName: "Alasri"}},
		Movie{ID: "2", Isbn: "45234", Title: "Movie2", Director: &Director{FirstName: "Doaa", LastName: "Massord"}},
		Movie{ID: "3", Isbn: "23453", Title: "Movie3", Director: &Director{FirstName: "Ali", LastName: "Ahemd"}},
	)

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movie", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at post 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
