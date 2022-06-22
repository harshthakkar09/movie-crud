package router

import (
	"movie-crud/src/authentication"
	"movie-crud/src/controllers"
	"movie-crud/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter() // defined new mux router

	r.HandleFunc("/login", authentication.Login).Methods("POST")
	r.HandleFunc("/signup", authentication.Signup).Methods("POST")
	MovieController := controllers.MovieController{}
	r.Handle("/movies", middlewares.Auth(http.HandlerFunc(MovieController.GetMovies))).Methods("GET")           // function handler to get all movies
	r.Handle("/movies/{id}", middlewares.Auth(http.HandlerFunc(MovieController.GetMovie))).Methods("GET")       // function handler to get a movie
	r.Handle("/movies", middlewares.Auth(http.HandlerFunc(MovieController.CreateMovie))).Methods("POST")        // function handler to create a movie
	r.Handle("/movies/{id}", middlewares.Auth(http.HandlerFunc(MovieController.UpdateMovie))).Methods("PUT")    // function handler to update a movie
	r.Handle("/movies/{id}", middlewares.Auth(http.HandlerFunc(MovieController.DeleteMovie))).Methods("DELETE") // function handler to delete a movie
	r.Handle("/movies/{id}/ratings", middlewares.Auth(http.HandlerFunc(MovieController.AddRatings))).Methods("POST")
	r.Handle("/movies/{id}/ratings", middlewares.Auth(http.HandlerFunc(MovieController.UpdateRatings))).Methods("PUT")
	r.Handle("/movies/{id}/ratings", middlewares.Auth(http.HandlerFunc(MovieController.DeleteRatings))).Methods("DELETE")

	CastController := controllers.CastController{}
	r.Handle("/cast", middlewares.Auth(http.HandlerFunc(CastController.GetCasts))).Methods("GET")           // function handler to get all casts
	r.Handle("/cast/{id}", middlewares.Auth(http.HandlerFunc(CastController.GetCast))).Methods("GET")       // function handler to get a cast
	r.Handle("/cast", middlewares.Auth(http.HandlerFunc(CastController.CreateCast))).Methods("POST")        // function handler to create a cast
	r.Handle("/cast/{id}", middlewares.Auth(http.HandlerFunc(CastController.UpdateCast))).Methods("PUT")    // function handler to update a cast
	r.Handle("/cast/{id}", middlewares.Auth(http.HandlerFunc(CastController.DeleteCast))).Methods("DELETE") // function handler to delete a cast

	return r
}
