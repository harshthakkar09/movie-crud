package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"movie-crud/src/models"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var (
	validGenre = []string{"thriller", "action", "horror", "fiction", "comedy"}
	movieMutex sync.Mutex
)

type MovieController struct{}

func checkStringInSlice(items []string, item string) bool {
	for _, cur := range items {
		if cur == item {
			return true
		}
	}
	return false
}

func ValidateMovieObject(movie *models.Movie) error {
	// checking whether movie title is present
	if movie.Title == "" {
		return fmt.Errorf("movie title is required")
	}

	// setting default value of genre
	if movie.Genre == "" {
		movie.Genre = "thriller"
	}

	// checking whether genre is valid
	flag := checkStringInSlice(validGenre, movie.Genre)
	if !flag {
		return fmt.Errorf("movie genre %s is not allowed", movie.Genre)
	}

	// checking whether all ratings have rater
	ratings := movie.Ratings
	for _, rating := range ratings {
		if rating.Rater == "" {
			return fmt.Errorf("rater is required")
		}
	}

	return nil
}

func castExists(castID string) bool {
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	json.Unmarshal(plan, &casts)

	// iterate through list of casts
	for _, item := range casts { // finding cast with given id
		if item.ID == castID {
			return true
		}
	}
	return false
}

func (m MovieController) CreateMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var movies []models.Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie) // decoding request body to Movie type of object
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	movie.ID = strconv.Itoa(rand.Intn(100000000)) // generating random id for Movie type of object

	err = ValidateMovieObject(&movie)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}

	// Validating cast-IDs
	for _, castID := range movie.CastIDs {
		if !castExists(castID) {
			http.Error(w, "cast with id="+castID+" does not exist", http.StatusBadRequest)
			return
		}
	}

	movies = append(movies, movie) // appending Movie type of object to movies array
	// writing updated movies array into json file
	file, _ := json.MarshalIndent(movies, "", " ")
	ioutil.WriteFile("./src/data/movies.json", file, 0644)
	json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
}

func (m MovieController) GetMovies(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var movies []models.Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// encoding and writing movies in json response
	json.NewEncoder(w).Encode(movies)
}

func (m MovieController) GetMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
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
	http.Error(w, "movie with given id not found", http.StatusNotFound)
}

func (m MovieController) UpdateMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var movies []models.Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// flag to check weather movie is found or not
	flag := false

	// iterate through list of movies
	for index, item := range movies {
		if item.ID == params["id"] { // finding movie with given id
			oldMovie := item
			var movie models.Movie
			err := json.NewDecoder(r.Body).Decode(&movie) // decoding request body to Movie type of object
			if err != nil {
				http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
				return
			}

			flag = true

			if oldMovie.Title != movie.Title {
				http.Error(w, "movie title cannot be changed", http.StatusBadRequest)
				return
			}

			err = ValidateMovieObject(&movie)
			if err != nil {
				http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
			}

			// Validating cast-IDs
			for _, castID := range movie.CastIDs {
				if !castExists(castID) {
					http.Error(w, "cast with id="+castID+" does not exist", http.StatusBadRequest)
					return
				}
			}

			movies = append(movies[:index], movies[index+1:]...) // updating movies array to delete a movie
			movie.ID = params["id"]                              // assigning id to Movie type of object
			movies = append(movies, movie)                       // appending Movie type of object to movies array

			// writing updated movies array into json file
			file, _ := json.MarshalIndent(movies, "", " ")
			ioutil.WriteFile("./src/data/movies.json", file, 0644)
			json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
			return
		}
	}

	// will throw error if no id found
	if !flag {
		http.Error(w, "movie with given id not found", http.StatusNotFound)
	}
}

func (m MovieController) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
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
			ioutil.WriteFile("./src/data/movies.json", file, 0644)
			// encoding and writing movies in json response
			json.NewEncoder(w).Encode(movies)
			return
		}
	}

	http.Error(w, "movie with given id not found", http.StatusNotFound)
}
