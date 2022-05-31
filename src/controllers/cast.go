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
	var castsMap map[string]interface{}
	err := json.Unmarshal(plan, &castsMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cast models.Cast
	err = json.NewDecoder(r.Body).Decode(&cast) // decoding request body to Cast object
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	cast.ID = uuid.New().String() // setting random id for cast object
	cast.ID = strings.ReplaceAll(cast.ID, "-", "")

	err = ValidateCastObject(&cast)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}

	castsMap[cast.ID] = cast // appending new cast object to castsMap

	// writing updated casts map into json file
	file, _ := json.MarshalIndent(castsMap, "", " ")
	ioutil.WriteFile("./src/data/casts.json", file, 0644)
	json.NewEncoder(w).Encode(cast) // encoding and writing cast in json response
}

func (c CastController) GetCasts(w http.ResponseWriter, r *http.Request) {
	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var castsMap map[string]interface{}
	err := json.Unmarshal(plan, &castsMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	casts := []models.Cast{}
	for _, val := range castsMap {
		str, err := json.Marshal(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var cast models.Cast
		err = json.Unmarshal(str, &cast)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		casts = append(casts, cast)
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
	var castsMap map[string]interface{}
	err := json.Unmarshal(plan, &castsMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// taking path parameters
	params := mux.Vars(r)

	castMap, ok := castsMap[params["id"]]

	if !ok {
		// will throw error if no id found
		http.Error(w, fmt.Sprintf("cast with id = %s not found", params["id"]), http.StatusNotFound)
	}

	str, err := json.Marshal(castMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cast models.Cast
	err = json.Unmarshal(str, &cast)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cast) // encoding and writing cast in json response
}

func (c CastController) UpdateCast(w http.ResponseWriter, r *http.Request) {

	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var castsMap map[string]interface{}
	err := json.Unmarshal(plan, &castsMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// taking path parameters
	params := mux.Vars(r)

	castMap, ok := castsMap[params["id"]]
	if !ok {
		// will throw error if no id found
		http.Error(w, fmt.Sprintf("cast with id = %s not found", params["id"]), http.StatusNotFound)
	}

	str, err := json.Marshal(castMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var oldCast models.Cast
	err = json.Unmarshal(str, &oldCast)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cast models.Cast
	err = json.NewDecoder(r.Body).Decode(&cast) // decoding request body to cast object
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
		return
	}

	// validating cast object
	err = ValidateCastObject(&cast)
	if err != nil {
		http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
	}

	delete(castsMap, oldCast.ID)  // removing old cast object from map
	cast.ID = params["id"]        // assigning id to cast object
	castsMap[params["id"]] = cast // adding new cast object to map

	// writing updated casts map into json file
	file, _ := json.MarshalIndent(castsMap, "", " ")
	ioutil.WriteFile("./src/data/casts.json", file, 0644)
	json.NewEncoder(w).Encode(cast) // encoding and writing cast in json response
}

func (c CastController) DeleteCast(w http.ResponseWriter, r *http.Request) {

	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var castsMap map[string]interface{}
	err := json.Unmarshal(plan, &castsMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// taking path parameters
	params := mux.Vars(r)

	_, ok := castsMap[params["id"]]
	if !ok {
		// will throw error if no id found
		http.Error(w, fmt.Sprintf("cast with id = %s not found", params["id"]), http.StatusNotFound)
	}
	delete(castsMap, params["id"]) // removing cast object from map

	// writing updated casts map into json file
	file, _ := json.MarshalIndent(castsMap, "", " ")
	ioutil.WriteFile("./src/data/casts.json", file, 0644)
	json.NewEncoder(w).Encode(fmt.Sprintf("cast with id = %s deleted successfully", params["id"])) // encoding and writing cast in json response
}
