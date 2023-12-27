package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json: "id"`
	Isbn     string    `json: "isbn"`
	Title    string    `json: "title"`
	Director *Director `json: "director`
}

type Director struct {
	Firstname string `json: "firstname"`
	Lastname  string `json: "lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/JSON")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/JSON")
	params := mux.Vars(r)
	for index, movie := range movies {
		if movie.ID == params["id"] {
			// Need to add ... at the end because
			// cannot use movies[index + 1:] (value of type []Movie) as Movie value in argument to append
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/JSON")
	params := mux.Vars(r)
	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/JSON")
	var movie Movie
	// When send as BODY not params, cannot use the following
	// params := mux.Vars(r)
	// movie := &Movie{ID: params["id"], Isbn: params["isbn"], Title: params["title"], Director: params&["director"]}
	_ = json.NewDecoder(r.Body).Decode(&movie)

	//randomly generate ID
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	// 1. set json content type
	// 2. params, find movie with the id
	// 3. delete it, add a new movie with details sent from the Request Body
	w.Header().Set("Content-Type", "application/JSON")
	params := mux.Vars(r)
	for index, movie := range movies {
		if movie.ID == params["id"] {
			// delete
			movies = append(movies[:index], movies[index+1:]...)
			var newMovie Movie
			_ = json.NewDecoder(r.Body).Decode(&newMovie)
			newMovie.ID = params["id"]
			movies = append(movies, newMovie)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "102464", Title: "Avengers", Director: &Director{Firstname: "Steven", Lastname: "Wong"}})
	movies = append(movies, Movie{ID: "2", Isbn: "205746", Title: "Ben10", Director: &Director{Firstname: "LeBron", Lastname: "James"}})
	movies = append(movies, Movie{ID: "3", Isbn: "365610", Title: "Cat", Director: &Director{Firstname: "Anthony", Lastname: "Davis"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Print("Starting server at PORT 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
