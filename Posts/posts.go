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

type Posts struct {
	gorm.Model
	Name  string
	Content string
	Image string
}

func allPosts(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var posts []Posts
	db.Find(&posts)
	json.NewEncoder(w).Encode(posts)
}


func createPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var posts Posts
	json.NewDecoder(r.Body).Decode(&posts)
	db.Create(&posts)
	json.NewEncoder(w).Encode(&posts)
}


//func newPosts(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("New User Endpoint Hit")
//
//	db, err := gorm.Open("sqlite3", "test.db")
//	if err != nil {
//		panic("failed to connect database")
//	}
//	defer db.Close()
//
//	vars := mux.Vars(r)
//	name := vars["name"]
//	content := vars["content"]
//	image := vars["image"]
//
//	db.Create(&Posts{Name: name, Content: content,Image:image})
//	fmt.Fprintf(w, "New User Successfully Created")
//}

func deletePosts(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	var posts Posts
	db.Where("name = ?", name).Find(&posts)
	db.Delete(&posts)

	fmt.Fprintf(w, "Successfully Deleted User")
}

func updatePosts(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]


	var posts Posts
	db.Where("name = ?", name).Find(&posts)

	posts.Name = name

	db.Save(&posts)
	fmt.Fprintf(w, "Successfully Updated User")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/posts", allPosts).Methods("GET")
	myRouter.HandleFunc("/posts/delete/{name}", deletePosts).Methods("DELETE")
	myRouter.HandleFunc("/posts/put/{name}/{content}/{image}", updatePosts).Methods("PUT")
	myRouter.HandleFunc("/posts/post", createPosts).Methods("POST")
	log.Fatal(http.ListenAndServe(":8010", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Posts{})
}

func main() {
	fmt.Println("Go ORM Tutorial")

	initialMigration()
	// Handle Subsequent requests
	handleRequests()
}
