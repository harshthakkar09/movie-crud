package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	//decode body for credentials
	var user Credentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, "Error in reading body")
		return
	}

	if len(user.Username) == 0 || len(user.Password) == 0 {
		WriteResponse(w, http.StatusBadRequest, "Username and password can not be empty")
		return
	}

	// Reading users.json to usersMap
	userJson, _ := ioutil.ReadFile("./src/data/users.json")
	var usersMap map[string]interface{}
	json.Unmarshal(userJson, &usersMap)

	_, ok := usersMap[user.Username]
	if ok {
		WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("User %s already exists", user.Username))
		return
	}

	//Generate hash for given password
	user.Password, err = GenerateHashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Append new user
	usersMap[user.Username] = user.Password
	file, _ := json.Marshal(usersMap)
	ioutil.WriteFile("./src/data/users.json", file, 0644)

	json.NewEncoder(w).Encode(fmt.Sprintf("User %s created successfully", user.Username))
}
