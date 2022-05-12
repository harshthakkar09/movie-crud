package router

import (
	"movie-crud/src/controllers"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter() // defined new mux router

	r.HandleFunc("/login", controllers.Login).Methods("POST")

	MovieController := controllers.MovieController{}
	r.HandleFunc("/movies", MovieController.GetMovies).Methods("GET")           // function handler to get all movies
	r.HandleFunc("/movies/{id}", MovieController.GetMovie).Methods("GET")       // function handler to get a movie
	r.HandleFunc("/movies", MovieController.CreateMovie).Methods("POST")        // function handler to create a movie
	r.HandleFunc("/movies/{id}", MovieController.UpdateMovie).Methods("PUT")    // function handler to update a movie
	r.HandleFunc("/movies/{id}", MovieController.DeleteMovie).Methods("DELETE") // function handler to delete a movie

	CastController := controllers.CastController{}
	r.HandleFunc("/cast", CastController.GetCasts).Methods("GET")           // function handler to get all casts
	r.HandleFunc("/cast/{id}", CastController.GetCast).Methods("GET")       // function handler to get a cast
	r.HandleFunc("/cast", CastController.CreateCast).Methods("POST")        // function handler to create a cast
	r.HandleFunc("/cast/{id}", CastController.UpdateCast).Methods("PUT")    // function handler to update a cast
	r.HandleFunc("/cast/{id}", CastController.DeleteCast).Methods("DELETE") // function handler to delete a cast

	return r
}
