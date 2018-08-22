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

type Events struct {
	gorm.Model
	Name  string
	Date string
	Meetup string
}

var events []Events

func allEvents(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()


	db.Find(&events)

	json.NewEncoder(w).Encode(events)
}

func newEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	date := vars["date"]
	meetup := vars["meetup"]

	db.Create(&Events{Name: name, Date: date,Meetup:meetup})
	fmt.Fprintf(w, "New User Successfully Created")
}

func deleteEvents(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	db.Where("name = ?", name).Find(&events)
	db.Delete(&events)

	fmt.Fprintf(w, "Successfully Deleted User")
}

func updateEvents(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	var event Events
	db.Where("name = ?", name).Find(&event)

	event.Name = name

	db.Save(&event)
	fmt.Fprintf(w, "Successfully Updated User")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/events", allEvents).Methods("GET")
	myRouter.HandleFunc("/event/delete/{name}", deleteEvents).Methods("DELETE")
	myRouter.HandleFunc("/event/put/{name}/{date}/{meetup}", updateEvents).Methods("PUT")
	myRouter.HandleFunc("/event/post/{name}/{date}/{meetup}", newEvents).Methods("POST")
	log.Fatal(http.ListenAndServe(":8020", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Events{})
}

func main() {
	fmt.Println("Go ORM Tutorial")

	initialMigration()
	// Handle Subsequent requests
	handleRequests()
}
