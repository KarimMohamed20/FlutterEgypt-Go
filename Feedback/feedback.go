package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Feedback struct {
	gorm.Model
	Name  string
	Content string
}

func allFeedback(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var feedback []Feedback
	db.Find(&feedback)
	json.NewEncoder(w).Encode(feedback)
}

func newFeedback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	content := vars["content"]

	fmt.Println(name)
	fmt.Println(content)

	db.Create(&Feedback{Name: name, Content: content})
	fmt.Fprintf(w, "New User Successfully Created")
}

func deleteFeedback(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var feedback Feedback
	db.Where("name = ?", name).Find(&feedback)
	db.Delete(&feedback)

	fmt.Fprintf(w, "Successfully Deleted User")
}

func updateFeedback(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	content := vars["content"]

	var feedback Feedback
	db.Where("name = ?", name).Find(&feedback)

	feedback.Content = content

	db.Save(&feedback)
	fmt.Fprintf(w, "Successfully Updated User")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/feedback", allFeedback).Methods("GET")
	myRouter.HandleFunc("/feedback/delete/{name}", deleteFeedback).Methods("DELETE")
	myRouter.HandleFunc("/feedback/put/{name}/{content}", updateFeedback).Methods("PUT")
	myRouter.HandleFunc("/feedback/post/{name}/{content}", newFeedback).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Feedback{})
}

func main() {
	fmt.Println("Go ORM Tutorial")

	initialMigration()
	// Handle Subsequent requests
	handleRequests()
}
