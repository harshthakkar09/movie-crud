package authorization

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWT Key used to create the signature
var jwtKey = []byte("cds/movie-provider")

func GetJWTKey() []byte {
	return jwtKey
}

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
		return
	}

	expectedPassword, ok := users[credentials.Username]

	// if password is not saved for given username
	// or password does not match
	if !ok || expectedPassword != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// expiration time : 5min
	expirationTime := time.Now().Add(5 * time.Minute)

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
}
