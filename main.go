package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	validGeneder = []string{"male", "female"}
	validGenre   = []string{"thriller", "action", "horror", "fiction", "comedy"}
)

// struct for movie object
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn,omitempty"`
	Title    string    `json:"title"`
	Director *Director `json:"director,omitempty"`
	CastIDs  []string  `json:"casts,omitempty"`
	Ratings  []Rating  `json:"ratings,omitempty"`
	Genre    string    `json:"genre"`
}

// struct for rating object
type Rating struct {
	Rater  string  `json:"rater"`
	Rating float32 `json:"rating"`
}

// struct for cast object
type Cast struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

// struct for director object
type Director struct {
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

func checkStringInSlice(items []string, item string) bool {
	for _, cur := range items {
		if cur == item {
			return true
		}
	}
	return false
}

func ValidateMovieObject(movie *Movie) error {
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

func ValidateCastObject(cast *Cast) error {
	// checking whether cast name is present
	if cast.Name == "" {
		return fmt.Errorf("cast name is required")
	}

	// setting default value of gender
	if cast.Gender == "" {
		cast.Gender = "male"
	}

	// check whether gender is valid
	flag := checkStringInSlice(validGeneder, cast.Gender)
	if !flag {
		return fmt.Errorf("cast gender %s is not allowed", cast.Gender)
	}

	return nil
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies.json")
	var movies []Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// encoding and writing movies in json response
	json.NewEncoder(w).Encode(movies)
}

func getCasts(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("casts.json")
	var casts []Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// encoding and writing casts in json response
	json.NewEncoder(w).Encode(casts)
}

func deleteCast(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("casts.json")
	var casts []Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// flag to check whether cast is found or not
	flag := false

	// iterate through list of casts
	for index, item := range casts {
		if item.ID == params["id"] { // finding cast with given id
			casts = append(casts[:index], casts[index+1:]...) // updating casts array to delete a cast
			flag = true                                       //updating flag value
			// writing updated casts array into json file
			file, _ := json.MarshalIndent(casts, "", " ")
			ioutil.WriteFile("casts.json", file, 0644)
			break
		}
	}

	// will throw error if no id found
	if !flag {
		http.Error(w, "cast with given id not found", http.StatusNotFound)
		return
	}

	// encoding and writing movies in json response
	json.NewEncoder(w).Encode(casts)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies.json")
	var movies []Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// flag to check whether movie is found or not
	flag := false

	// iterate through list of movies
	for index, item := range movies {
		if item.ID == params["id"] { // finding movie with given id
			movies = append(movies[:index], movies[index+1:]...) // updating movies array to delete a movie
			flag = true                                          //updating flag value
			// writing updated movies array into json file
			file, _ := json.MarshalIndent(movies, "", " ")
			ioutil.WriteFile("movies.json", file, 0644)
			break
		}
	}

	// will throw error if no id found
	if !flag {
		http.Error(w, "movie with given id not found", http.StatusNotFound)
		return
	}

	// encoding and writing movies in json response
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies.json")
	var movies []Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// flag to check whether movie is found or not
	flag := false

	// iterate through list of movies
	for _, item := range movies { // finding movie with given id
		if item.ID == params["id"] {
			flag = true //updating flag value
			// writing updated movies array into json file
			json.NewEncoder(w).Encode(item) // encoding and writing movie in json response
			return
		}
	}

	// will throw error if no id found
	if !flag {
		http.Error(w, "movie with given id not found", http.StatusNotFound)
	}
}

func getCast(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("casts.json")
	var casts []Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// flag to check whether cast is found or not
	flag := false

	// iterate through list of casts
	for _, item := range casts { // finding cast with given id
		if item.ID == params["id"] {
			flag = true                     //updating flag value
			json.NewEncoder(w).Encode(item) // encoding and writing cast in json response
			return
		}
	}

	// will throw error if no id found
	if !flag {
		http.Error(w, "cast with given id not found", http.StatusNotFound)
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies.json")
	var movies []Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	var movie Movie
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

	movies = append(movies, movie) // appending Movie type of object to movies array
	// writing updated movies array into json file
	file, _ := json.MarshalIndent(movies, "", " ")
	ioutil.WriteFile("movies.json", file, 0644)
	json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
}

func createCast(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("casts.json")
	var casts []Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	var cast Cast

	err := json.NewDecoder(r.Body).Decode(&cast) // decoding request body to Cast type of object
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}

	cast.ID = strconv.Itoa(rand.Intn(100000000)) // generating random id for Cast type of object

	err = ValidateCastObject(&cast)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}

	casts = append(casts, cast) // appending Movie type of object to casts array
	// writing updated casts array into json file
	file, _ := json.MarshalIndent(casts, "", " ")
	ioutil.WriteFile("casts.json", file, 0644)
	json.NewEncoder(w).Encode(cast) // encoding and writing movie in json response
}

func updateCast(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("casts.json")
	var casts []Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// flag to check whether cast is found or not
	flag := false

	// iterate through list of casts
	for index, item := range casts {
		if item.ID == params["id"] { // finding cast with given id
			casts = append(casts[:index], casts[index+1:]...) // updating casts array to delete a cast
			castOld := item                                   //will store old cast value for comperision purpose
			var cast Cast
			err := json.NewDecoder(r.Body).Decode(&cast) // decoding request body to Cast type of object
			if err != nil {
				http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
				return
			}

			flag = true

			// will not allow user to change cast name
			if castOld.Name != cast.Name {
				http.Error(w, "cast name cannot be changed", http.StatusBadRequest)
				return
			}

			cast.ID = params["id"] // assigning id to Cast type of object

			err = ValidateCastObject(&cast)
			if err != nil {
				http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
				return
			}

			casts = append(casts, cast) // appending Cast type of object to Casts array
			// writing updated casts array into json file
			file, _ := json.MarshalIndent(casts, "", " ")
			ioutil.WriteFile("casts.json", file, 0644)
			json.NewEncoder(w).Encode(cast) // encoding and writing cast in json response
			return
		}
	}

	// will throw error if no id found
	if !flag {
		http.Error(w, "cast with given id not found", http.StatusNotFound)
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("movies.json")
	var movies []Movie
	json.Unmarshal(plan, &movies)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// flag to check whether movie is found or not
	flag := false

	// iterate through list of movies
	for index, item := range movies {
		if item.ID == params["id"] { // finding movie with given id
			movies = append(movies[:index], movies[index+1:]...) // updating movies array to delete a movie
			oldMovie := item
			var movie Movie
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

			movie.ID = params["id"]        // assigning id to Movie type of object
			movies = append(movies, movie) // appending Movie type of object to movies array

			err = ValidateMovieObject(&movie)
			if err != nil {
				http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
			}

			// writing updated movies array into json file
			file, _ := json.MarshalIndent(movies, "", " ")
			ioutil.WriteFile("movies.json", file, 0644)
			json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
			return
		}
	}

	// will throw error if no id found
	if !flag {
		http.Error(w, "movie with given id not found", http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter() // defined new mux router

	r.HandleFunc("/movies", getMovies).Methods("GET")           // function handler to get all movies
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")       // function handler to get a movie
	r.HandleFunc("/movies", createMovie).Methods("POST")        // function handler to create a movie
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")    // function handler to update a movie
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE") // function handler to delete a movie

	r.HandleFunc("/cast", getCasts).Methods("GET")           // function handler to get all casts
	r.HandleFunc("/cast/{id}", getCast).Methods("GET")       // function handler to get a cast
	r.HandleFunc("/cast", createCast).Methods("POST")        // function handler to create a cast
	r.HandleFunc("/cast/{id}", updateCast).Methods("PUT")    // function handler to update a cast
	r.HandleFunc("/cast/{id}", deleteCast).Methods("DELETE") // function handler to delete a cast

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r)) //starting server on PORT 8000

}
