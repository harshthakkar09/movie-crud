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

type CastController struct{}

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

func (c CastController) DeleteCast(w http.ResponseWriter, r *http.Request) {
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

func (c CastController) GetCast(w http.ResponseWriter, r *http.Request) {
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

func (c CastController) CreateCast(w http.ResponseWriter, r *http.Request) {
	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var casts []models.Cast
	json.Unmarshal(plan, &casts)

	// defining header's content-type
	w.Header().Set("Content-Type", "application/json")

	var cast models.Cast
	_ = json.NewDecoder(r.Body).Decode(&cast)    // decoding request body to Cast type of object
	cast.ID = strconv.Itoa(rand.Intn(100000000)) // generating random id for Cast type of object
	casts = append(casts, cast)                  // appending Movie type of object to casts array
	// writing updated casts array into json file
	file, _ := json.MarshalIndent(casts, "", " ")
	ioutil.WriteFile("./src/data/casts.json", file, 0644)
	json.NewEncoder(w).Encode(cast) // encoding and writing movie in json response
}

func (c CastController) UpdateCast(w http.ResponseWriter, r *http.Request) {
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

			var cast models.Cast
			_ = json.NewDecoder(r.Body).Decode(&cast) // decoding request body to Cast type of object
			cast.ID = params["id"]                    // assigning id to Cast type of object
			casts = append(casts, cast)               // appending Cast type of object to Casts array
			// writing updated casts array into json file
			file, _ := json.MarshalIndent(casts, "", " ")
			ioutil.WriteFile("./src/data/casts.json", file, 0644)
			json.NewEncoder(w).Encode(cast) // encoding and writing cast in json response
			return
		}
	}
}
