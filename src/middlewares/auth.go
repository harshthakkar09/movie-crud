package middlewares

import (
	"net/http"

	"movie-crud/src/authentication"

	"github.com/golang-jwt/jwt"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				// cookie is not set
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// other errors
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// extracting JWT from cookie
		tokenString := c.Value

		claims := &authentication.Claims{}
		// parsing JWT string and storing the result in claims
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return authentication.GetJWTKey(), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				// invalid signature
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// defining header's content-type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
