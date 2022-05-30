package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	// decoding json body into credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(credentials.Username) == 0 || len(credentials.Password) == 0 {
		WriteResponse(w, http.StatusBadRequest, "Username and password can not be empty")
		return
	}

	// Reading users.json to usersMap
	userJson, _ := ioutil.ReadFile("./src/data/users.json")
	var usersMap map[string]interface{}
	json.Unmarshal(userJson, &usersMap)

	// getting hashed password from usersMap
	hashPassword, ok := usersMap[credentials.Username]
	if !ok {
		WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("User %s does not exist", credentials.Username))
		return
	}

	password := credentials.Password

	// comparing credentials
	if !CheckPasswordHash(password, hashPassword.(string)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// expiration time : 1 hour
	expirationTime := time.Now().Add(1 * time.Hour)

	// creating JWT claim using username and expiration time
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// creating and signing token for defined claims using HS256 method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Creating JWT string
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		// error on creating JWT
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// setting cookie for "token"
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode("Login Success! Token Expires within 1 Hour.")
}
