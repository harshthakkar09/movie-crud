package authentication

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	//decode body for credentials
	var user Credentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		err := "Error in reading body"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	//Reading users.json for user existence
	var usersArray []Credentials
	userjson, _ := ioutil.ReadFile("./src/data/users.json")
	json.Unmarshal(userjson, &usersArray)

	for _, v := range usersArray {
		if v.Username == user.Username {
			err := "Username already in use"
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err)
			return
		}
	}

	//Generate hash for given password
	user.Password, err = GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	//Append new user
	usersArray = append(usersArray, user)
	file, _ := json.MarshalIndent(usersArray, "", " ")
	ioutil.WriteFile("./src/data/users.json", file, 0644)
	json.NewEncoder(w).Encode(user)
}
