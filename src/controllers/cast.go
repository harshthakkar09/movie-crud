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
	validGender = []string{"male", "female"}
	castMutex   sync.Mutex
)

type CastController struct{}

func ValidateCastObject(cast *models.Cast) error {
	// checking whether cast name is present
	if cast.Name == "" {
		return fmt.Errorf("cast name is required")
	}

	// setting default value of gender
	if cast.Gender == "" {
		cast.Gender = "male"
	}

	// check whether gender is valid
	flag := checkStringInSlice(validGender, cast.Gender)
	if !flag {
		return fmt.Errorf("cast gender %s is not allowed", cast.Gender)
	}

	return nil
}

func (c CastController) CreateCast(w http.ResponseWriter, r *http.Request) {

	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	var cast models.Cast

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
	ioutil.WriteFile("./src/data/casts.json", file, 0644)
	json.NewEncoder(w).Encode(cast) // encoding and writing movie in json response
}

func (c CastController) GetCasts(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// encoding and writing casts in json response
	json.NewEncoder(w).Encode(casts)
}

func (c CastController) GetCast(w http.ResponseWriter, r *http.Request) {

	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)

	// iterate through list of casts
	for _, item := range casts { // finding cast with given id
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item) // encoding and writing cast in json response
			return
		}
	}
}

func (c CastController) UpdateCast(w http.ResponseWriter, r *http.Request) {

	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
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
			castOld := item //will store old cast value for comperision purpose
			var cast models.Cast
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

			casts = append(casts[:index], casts[index+1:]...) // updating casts array to delete a cast
			casts = append(casts, cast)                       // appending Cast type of object to Casts array
			// writing updated casts array into json file
			file, _ := json.MarshalIndent(casts, "", " ")
			ioutil.WriteFile("./src/data/casts.json", file, 0644)
			json.NewEncoder(w).Encode(cast) // encoding and writing cast in json response
			return
		}
	}

	// will throw error if no id found
	if !flag {
		http.Error(w, "cast with given id not found", http.StatusNotFound)
	}
}

func (c CastController) DeleteCast(w http.ResponseWriter, r *http.Request) {

	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	// taking path parameters
	params := mux.Vars(r)
	// iterate through list of casts
	for index, item := range casts {
		if item.ID == params["id"] { // finding cast with given id
			casts = append(casts[:index], casts[index+1:]...) // updating casts array to delete a cast
			// writing updated casts array into json file
			file, _ := json.MarshalIndent(casts, "", " ")
			ioutil.WriteFile("./src/data/casts.json", file, 0644)
			break
		}
	}
	// encoding and writing movies in json response
	json.NewEncoder(w).Encode(casts)
}
