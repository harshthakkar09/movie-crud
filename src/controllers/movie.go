package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"movie-crud/src/models"
	"net/http"
	"strconv"
	"strings"
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
		return fmt.Errorf("genre value %s is not allowed", movie.Genre)
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

func castExists(CastIDs []string) ([]string, bool, error) {
	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	err := json.Unmarshal(plan, &casts)
	if err != nil {
		return nil, false, err
	}
	// iterate through list of casts
	notPresentIDs := []string{}
	for _, castid := range CastIDs { // finding cast with castid
		found := false
		for _, cast := range casts {
			if cast.ID == castid {
				found = true
				break
			}
		}
		if !found {
			notPresentIDs = append(notPresentIDs, castid)
		}
	}
	if len(notPresentIDs) == 0 {
		return nil, true, nil
	}
	return notPresentIDs, false, nil
}

func (m MovieController) CreateMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var movies []models.Movie
	err := json.Unmarshal(plan, &movies)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	var movie models.Movie
	err = json.NewDecoder(r.Body).Decode(&movie) // decoding request body to Movie type of object
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	id := strconv.Itoa(rand.Intn(100000000))
	for {
		isValidId := true
		for _, v := range movies {
			if v.ID == id {
				isValidId = false
			}
		}
		if isValidId {
			break
		}
		id = strconv.Itoa(rand.Intn(100000000))
	}
	movie.ID = id // setting random id for Movie type of object

	err = ValidateMovieObject(&movie)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}

	// Validating cast-IDs
	notPresentIDs, res, err := castExists(movie.CastIDs)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}
	if !res {
		http.Error(w, "casts with id=["+strings.Join(notPresentIDs, ", ")+"] does not exist", http.StatusBadRequest)
		return
	}

	movies = append(movies, movie) // appending Movie type of object to movies array
	// writing updated movies array into json file
	file, _ := json.MarshalIndent(movies, "", " ")
	ioutil.WriteFile("./src/data/movies.json", file, 0644)
	json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
}

func (m MovieController) GetMovies(w http.ResponseWriter, r *http.Request) {
	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var movies []models.Movie
	err := json.Unmarshal(plan, &movies)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

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
	err := json.Unmarshal(plan, &movies)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	// taking path parameters
	params := mux.Vars(r)

	// iterate through list of movies
	for _, item := range movies { // finding movie with given id
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item) // encoding and writing movie in json response
			return
		}
	}

	// will throw error if no id found
	http.Error(w, fmt.Sprintf("movie with id = %s not found", params["id"]), http.StatusNotFound)
}

func (m MovieController) UpdateMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var movies []models.Movie
	err := json.Unmarshal(plan, &movies)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	// taking path parameters
	params := mux.Vars(r)

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

			if oldMovie.Title != movie.Title {
				http.Error(w, "movie title cannot be changed", http.StatusBadRequest)
				return
			}

			err = ValidateMovieObject(&movie)
			if err != nil {
				http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
			}

			// Validating cast-IDs
			notPresentIDs, res, err := castExists(movie.CastIDs)
			if err != nil {
				http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
				return
			}
			if !res {
				http.Error(w, "casts with id=["+strings.Join(notPresentIDs, ", ")+"] does not exist", http.StatusBadRequest)
				return
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
	http.Error(w, fmt.Sprintf("movie with id = %s not found", params["id"]), http.StatusNotFound)
}

func (m MovieController) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var movies []models.Movie
	err := json.Unmarshal(plan, &movies)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

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

	// will throw error if no id found
	http.Error(w, fmt.Sprintf("movie with id = %s not found", params["id"]), http.StatusNotFound)
}
