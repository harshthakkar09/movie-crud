package main

import (
	"fmt"
	"log"
	"net/http"

	"movie-crud/src/router"
)

func main() {
	r := router.RegisterRoutes()

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r)) //starting server on PORT 8000
}
