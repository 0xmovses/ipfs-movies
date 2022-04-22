package routes

import (
	"movie-mux/pkg/controllers"

	"github.com/gorilla/mux"
)

//RegisterMovieRoutes registers the movie routes
var RegisterMovieRoutes = func(router *mux.Router) {
	router.HandleFunc("/movies/", controllers.GetMovies).Methods("GET")
	router.HandleFunc("/movies/", controllers.CreateMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", controllers.GetMovieByID).Methods("GET")
	router.HandleFunc("/movies/{id}", controllers.UpdateMovie).Methods("PUT")
}
