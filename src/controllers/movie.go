package controllers

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"movie-crud/src/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MovieController struct{}

func (m MovieController) CreateMovie(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies-crud/src/data/movies.json")
	var movies []models.Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	var movie models.Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)    // decoding request body to Movie type of object
	movie.ID = strconv.Itoa(rand.Intn(100000000)) // generating random id for Movie type of object
	movies = append(movies, movie)                // appending Movie type of object to movies array
	// writing updated movies array into json file
	file, _ := json.MarshalIndent(movies, "", " ")
	ioutil.WriteFile("movies.json", file, 0644)
	json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
}

func (m MovieController) GetMovies(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies.json")
	var movies []models.Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// encoding and writing movies in json response
	json.NewEncoder(w).Encode(movies)
}

func (m MovieController) GetMovie(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies.json")
	var movies []models.Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// iterate through list of movies
	for _, item := range movies { // finding movie with given id
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item) // encoding and writing movie in json response
			return
		}
	}
}

func (m MovieController) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies.json")
	var movies []models.Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// iterate through list of movies
	for index, item := range movies {
		if item.ID == params["id"] { // finding movie with given id
			movies = append(movies[:index], movies[index+1:]...) // updating movies array to delete a movie

			var movie models.Movie
			_ = json.NewDecoder(r.Body).Decode(&movie) // decoding request body to Movie type of object
			movie.ID = params["id"]                    // assigning id to Movie type of object
			movies = append(movies, movie)             // appending Movie type of object to movies array
			// writing updated movies array into json file
			file, _ := json.MarshalIndent(movies, "", " ")
			ioutil.WriteFile("movies.json", file, 0644)
			json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
			return
		}
	}
}

func (m MovieController) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies.json")
	var movies []models.Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// iterate through list of movies
	for index, item := range movies {
		if item.ID == params["id"] { // finding movie with given id
			movies = append(movies[:index], movies[index+1:]...) // updating movies array to delete a movie
			// writing updated movies array into json file
			file, _ := json.MarshalIndent(movies, "", " ")
			ioutil.WriteFile("movies.json", file, 0644)
			break
		}
	}
	// encoding and writing movies in json response
	json.NewEncoder(w).Encode(movies)
}
