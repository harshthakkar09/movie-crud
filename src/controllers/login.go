package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
)

// JWT Key used to create the signature
var jwtKey = []byte("cds/movie-provider")

var users = map[string]string{
	"admin": "Crest@123",
	"guest": "Guest@123",
}

// Structure used to extract username and password from request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// A structure which will be converted to JWT
// Username and expiration time(from standard claims)
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	// decoding json body into credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		// decode failed, wrong json body
		w.WriteHeader(http.StatusBadRequest)
	}
}
