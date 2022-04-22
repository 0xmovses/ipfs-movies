package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"movie-mux/pkg/models"

	"movie-mux/pkg/utils"
)

var NewMovie models.Movie

func GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := models.GetAllMovies()
	if err != nil {

	}
	res, _ := json.Marshal(movies)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateMovie")
	Movie := &models.Movie{}
	utils.ParseBody(r, Movie)
	fmt.Println("parsed body")
	m := Movie.Create()
	fmt.Println("created")
	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
