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
		return fmt.Errorf("gender value %s is not allowed", cast.Gender)
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
	err := json.Unmarshal(plan, &casts)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	var cast models.Cast

	err = json.NewDecoder(r.Body).Decode(&cast) // decoding request body to Cast type of object
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}
	id := strconv.Itoa(rand.Intn(100000000))
	for {
		isValidId := true
		for _, v := range casts {
			if v.ID == id {
				isValidId = false
			}
		}
		if isValidId {
			break
		}
		id = strconv.Itoa(rand.Intn(100000000))
	}
	cast.ID = id // generating random id for Cast type of object
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
	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	err := json.Unmarshal(plan, &casts)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

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
	err := json.Unmarshal(plan, &casts)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	// taking path parameters
	params := mux.Vars(r)

	// iterate through list of casts
	for _, item := range casts { // finding cast with given id
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item) // encoding and writing cast in json response
			return
		}
	}

	// will throw error if no id found
	http.Error(w, fmt.Sprintf("cast with id = %s not found", params["id"]), http.StatusNotFound)
}

func (c CastController) UpdateCast(w http.ResponseWriter, r *http.Request) {

	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	err := json.Unmarshal(plan, &casts)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	// taking path parameters
	params := mux.Vars(r)

	// iterate through list of casts
	for index, item := range casts {
		if item.ID == params["id"] { // finding cast with given id
			castOld := item //will store old cast value for comparison purpose
			var cast models.Cast
			err := json.NewDecoder(r.Body).Decode(&cast) // decoding request body to Cast type of object
			if err != nil {
				http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
				return
			}

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
	http.Error(w, fmt.Sprintf("cast with id = %s not found", params["id"]), http.StatusNotFound)

}

func (c CastController) DeleteCast(w http.ResponseWriter, r *http.Request) {

	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	err := json.Unmarshal(plan, &casts)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	// taking path parameters
	params := mux.Vars(r)
	// iterate through list of casts
	for index, item := range casts {
		if item.ID == params["id"] { // finding cast with given id
			casts = append(casts[:index], casts[index+1:]...) // updating casts array to delete a cast
			// writing updated casts array into json file
			file, _ := json.MarshalIndent(casts, "", " ")
			ioutil.WriteFile("./src/data/casts.json", file, 0644)
			// encoding and writing movies in json response
			json.NewEncoder(w).Encode(casts)
			return
		}
	}

	// will throw error if no id found
	http.Error(w, fmt.Sprintf("cast with id = %s not found", params["id"]), http.StatusNotFound)
}
