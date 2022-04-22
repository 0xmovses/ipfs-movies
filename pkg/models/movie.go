package models

import (
	"fmt"
	"log"
	"movie-mux/pkg/config"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	ID          int64  `gorm:"" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Movie{})
}

func GetAllMovies() ([]Movie, error) {
	var movies []Movie
	err := db.Find(&movies).Error
	return movies, err
}

func (m *Movie) Create() *Movie {
	fmt.Println("Creating movie")
	result := db.Create(&m)
	//log the result
	log.Println(result)
	return m
}
