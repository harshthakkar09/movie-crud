package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"movie-crud/src/models"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	validGenre = []string{"thriller", "action", "horror", "fiction", "comedy"}
	movieMutex sync.Mutex
)

type MovieController struct{}

func (m MovieController) CreateMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var moviesMap map[string]interface{}
	err := json.Unmarshal(plan, &moviesMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var movie models.Movie
	err = json.NewDecoder(r.Body).Decode(&movie) // decoding request body to Movie object
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	for _, val := range moviesMap {
		str, err := json.Marshal(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var movieRemote models.Movie
		err = json.Unmarshal(str, &movieRemote)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// checking for duplicate title
		if movie.Title == movieRemote.Title {
			http.Error(w, fmt.Sprintf("movie with title %s already exists", movieRemote.Title), http.StatusBadRequest)
			return
		}
	}

	movie.ID = uuid.New().String() // setting random id for Movie object
	movie.ID = strings.ReplaceAll(movie.ID, "-", "")

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

	moviesMap[movie.ID] = movie // appending new movie object to moviesMap

	// writing updated movies map into json file
	file, _ := json.MarshalIndent(moviesMap, "", " ")
	ioutil.WriteFile("./src/data/movies.json", file, 0644)
	json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
}

func (m MovieController) GetMovies(w http.ResponseWriter, r *http.Request) {
	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var moviesMap map[string]interface{}
	err := json.Unmarshal(plan, &moviesMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	movies := []models.Movie{}
	for _, val := range moviesMap {
		str, err := json.Marshal(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var movie models.Movie
		err = json.Unmarshal(str, &movie)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		movies = append(movies, movie)
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
	var moviesMap map[string]interface{}
	err := json.Unmarshal(plan, &moviesMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// taking path parameters
	params := mux.Vars(r)

	movieMap, ok := moviesMap[params["id"]]

	if !ok {
		// will throw error if no id found
		http.Error(w, fmt.Sprintf("movie with id = %s not found", params["id"]), http.StatusNotFound)
	}

	str, err := json.Marshal(movieMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var movie models.Movie
	err = json.Unmarshal(str, &movie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
}

func (m MovieController) UpdateMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var moviesMap map[string]interface{}
	err := json.Unmarshal(plan, &moviesMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// taking path parameters
	params := mux.Vars(r)

	movieMap, ok := moviesMap[params["id"]]
	if !ok {
		// will throw error if no id found
		http.Error(w, fmt.Sprintf("movie with id = %s not found", params["id"]), http.StatusNotFound)
	}

	str, err := json.Marshal(movieMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var oldMovie models.Movie
	err = json.Unmarshal(str, &oldMovie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var movie models.Movie
	err = json.NewDecoder(r.Body).Decode(&movie) // decoding request body to Movie object
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}

	if oldMovie.Title != movie.Title {
		http.Error(w, "movie title cannot be changed", http.StatusBadRequest)
		return
	}

	// validating movie object
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

	delete(moviesMap, oldMovie.ID)  // removing old movie object from map
	movie.ID = params["id"]         // assigning id to Movie object
	moviesMap[params["id"]] = movie // adding new movie object to map

	// writing updated movies map into json file
	file, _ := json.MarshalIndent(moviesMap, "", " ")
	ioutil.WriteFile("./src/data/movies.json", file, 0644)
	json.NewEncoder(w).Encode(movie) // encoding and writing movie in json response
}

func (m MovieController) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var moviesMap map[string]interface{}
	err := json.Unmarshal(plan, &moviesMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// taking path parameters
	params := mux.Vars(r)

	_, ok := moviesMap[params["id"]]
	if !ok {
		// will throw error if no id found
		http.Error(w, fmt.Sprintf("movie with id = %s not found", params["id"]), http.StatusNotFound)
	}
	delete(moviesMap, params["id"]) // removing movie object from map

	// writing updated movies map into json file
	file, _ := json.MarshalIndent(moviesMap, "", " ")
	ioutil.WriteFile("./src/data/movies.json", file, 0644)
	json.NewEncoder(w).Encode(fmt.Sprintf("Movie with id = %s deleted successfully", params["id"])) // encoding and writing movie in json response
}

func (m MovieController) AddRatings(w http.ResponseWriter, r *http.Request) {
	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var moviesMap map[string]interface{}
	err := json.Unmarshal(plan, &moviesMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// var movie models.Movie

	params := mux.Vars(r)

	mov, ok := moviesMap[params["id"]]
	if !ok {
		http.Error(w, fmt.Sprintf("movie with id = %s not found", params["id"]), http.StatusNotFound)
	}

	movieMap := mov.(map[string]interface{})

	var ratingData models.Rating
	err = json.NewDecoder(r.Body).Decode(&ratingData)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}
	if ratsObj, exists := movieMap["ratings"]; exists {
		var ratings []models.Rating
		ratingBytes, _ := json.Marshal(ratsObj)
		json.Unmarshal(ratingBytes, &ratings)

		for _, re := range ratings {
			if re.Rater == ratingData.Rater {
				http.Error(w, fmt.Sprintf("Rater Already Exists...."), http.StatusOK)
				return
			}
		}
		ratings = append(ratings, ratingData)
		// fmt.Println(ratings)
		// r, _ := json.Marshal(ratings)

		movieMap["ratings"] = ratings

		var movie models.Movie
		movieBytes, _ := json.Marshal(movieMap)
		json.Unmarshal(movieBytes, &movie)
		fmt.Println(movie)
		moviesMap[params["id"]] = movie

		// writing updated movies map into json file
		file, _ := json.MarshalIndent(moviesMap, "", " ")
		ioutil.WriteFile("./src/data/movies.json", file, 0644)
		err = json.NewEncoder(w).Encode(movieMap)
	} else {

		var ratings []models.Rating
		ratings = append(ratings, ratingData)
		movieMap["ratings"] = ratings
		var movie models.Movie
		movieBytes, _ := json.Marshal(movieMap)
		json.Unmarshal(movieBytes, &movie)
		moviesMap[params["id"]] = movie

		// writing updated movies map into json file
		file, _ := json.MarshalIndent(moviesMap, "", " ")
		ioutil.WriteFile("./src/data/movies.json", file, 0644)
		err = json.NewEncoder(w).Encode(movieMap)
	}

}

func (m MovieController) UpdateRatings(w http.ResponseWriter, r *http.Request) {
	//mutex
	movieMutex.Lock()
	defer movieMutex.Unlock()

	plan, _ := ioutil.ReadFile("./src/data/movies.json")
	var moviesMap map[string]interface{}
	err := json.Unmarshal(plan, &moviesMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)

	mov, ok := moviesMap[params["id"]]
	if !ok {
		http.Error(w, fmt.Sprintf("movie with id = %s not found", params["id"]), http.StatusNotFound)
	}

	movieMap := mov.(map[string]interface{})

	var ratingData models.Rating
	err = json.NewDecoder(r.Body).Decode(&ratingData)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	if ratsObj, exists := movieMap["ratings"]; exists {
		var ratings []models.Rating
		ratingBytes, _ := json.Marshal(ratsObj)
		json.Unmarshal(ratingBytes, &ratings)

		for i := range ratings {
			if ratings[i].Rater == ratingData.Rater {
				ratings[i].Rating = ratingData.Rating

				movieMap["ratings"] = ratings
				var movie models.Movie
				movieBytes, _ := json.Marshal(movieMap)
				json.Unmarshal(movieBytes, &movie)
				fmt.Println(movie)
				moviesMap[params["id"]] = movie

			} else {
				http.Error(w, fmt.Sprintf("Rater Doesn't Exists...."), http.StatusOK)
				return
			}
		}
		// writing updated movies map into json file
		file, _ := json.MarshalIndent(moviesMap, "", " ")
		ioutil.WriteFile("./src/data/movies.json", file, 0644)
		err = json.NewEncoder(w).Encode(movieMap)
	}
}
