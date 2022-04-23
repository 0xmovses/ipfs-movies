package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"movie-mux/pkg/models"
	"movie-mux/pkg/shared"

	"movie-mux/pkg/utils"

	"github.com/gorilla/mux"
)

//NewMovie is of type Movie
var NewMovie models.Movie

//GetMovies returns all movies
func GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := models.GetAllMovies()
	if err != nil {

	}
	res, _ := json.Marshal(movies)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

//CreateMovie creates a new movie
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	Movie := &models.Movie{}
	utils.ParseBody(r, Movie)
	m := Movie.Create()

	//create a new Ipfs Client
	ipfsClient := shared.NewIpfsClient(
		os.Getenv("PINATA_API_KEY"),
		os.Getenv("PINATA_API_SECRET"))

	w.Header().Set("Content-Type", "application/json")

	_, err := ipfsClient.PinJSON(m)
	if err != nil {
		fmt.Println("error pinning json to ipfs", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
	}

	res, _ := json.Marshal(m)

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

//GetMovieByID returns a movie by id
func GetMovieByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ID, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, _ := models.GetMovieById(ID)
	res, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

//UpdateMovie updates a movie by id
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	var updateMovie = &models.Movie{}
	utils.ParseBody(r, updateMovie)
	vars := mux.Vars(r)
	id := vars["id"]
	ID, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, db := models.GetMovieById(ID)
	if updateMovie.Title != "" {
		m.Title = updateMovie.Title
	}
	if updateMovie.Description != "" {
		m.Description = updateMovie.Description
	}
	db.Save(&m)
	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
