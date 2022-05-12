package authentication

import (
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// JWT Key used to create the signature
var jwtKey = []byte("cds/movie-provider")

func GetJWTKey() []byte {
	return jwtKey
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

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
