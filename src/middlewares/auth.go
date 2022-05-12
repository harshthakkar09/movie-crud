package middlewares

import (
	"fmt"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middle ware...")
		next.ServeHTTP(w, r)
	})
}
