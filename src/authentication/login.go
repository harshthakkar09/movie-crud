package authentication

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	//Reading users.json to Array
	var usersArray []Credentials
	userjson, _ := ioutil.ReadFile("./src/data/users.json")
	json.Unmarshal(userjson, &usersArray)

	// decoding json body into credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Get hashpassword for given username
	var hashPassword string
	for _, v := range usersArray {
		if v.Username == credentials.Username {
			hashPassword = v.Password
		}
	}

	password := credentials.Password

	if !CheckPasswordHash(password, hashPassword) {
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

	loginmsg := "login Sucess Token Expires within 5 minutes."
	json.NewEncoder(w).Encode(loginmsg)
}
