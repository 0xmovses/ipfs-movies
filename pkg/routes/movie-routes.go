package routes

import (
	"movie-mux/pkg/controllers"

	"github.com/gorilla/mux"
)

var RegisterMovieRoutes = func(router *mux.Router) {
	router.HandleFunc("/movies/", controllers.GetMovies).Methods("GET")
	router.HandleFunc("/movies/", controllers.CreateMovie).Methods("POST")
}
