/*  No DB is used for this practice
    gorilla/mux used for routing
    Testing is done using Postman */
    
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// json tags are useful for unmarshalling request bodies
type randomInfoAboutMe struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Content string `json:"content"`
	Desc    string `json:"desc"`
}

// Entries populated in main, used as a simulation for DB
var Entries []randomInfoAboutMe

// homePage welcomes you as a sign of Turkish Hospitality
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the realm of HTTP Wizards!")
}

// returnSingleEntry returns requested id's entry as JSON if it exists
func returnSingleEntry(w http.ResponseWriter, r *http.Request) {
	// Get id to key variable
	vars := mux.Vars(r)
	key := vars["id"]

	// Loop over all of our Entries return the entry encoded as JSON
	for _, entry := range Entries {
		if entry.Id == key {
			json.NewEncoder(w).Encode(entry)
		}
	}
}

// returnAllEntries returns all entries, if there is any
func returnAllEntries(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Entries)
}

// createNewEntry handles POST request with a body by appending it to global sturctu Entries
func createNewEntry(w http.ResponseWriter, r *http.Request) {
	// Get the body of our POST request
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newInfo randomInfoAboutMe

	// JSON to randomInfoAboutMe(My Struct)
	json.Unmarshal(reqBody, &newInfo)

	// Update our global Entries array with newly added entry
	Entries = append(Entries, newInfo)

	// This line prints the body of the request it received to inform the client (double-checking)
	json.NewEncoder(w).Encode(newInfo)
}

// updateEntry updates only the pointed(by id) struct with the given body information
func updateEntry(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var updateInfo randomInfoAboutMe
	json.Unmarshal(reqBody, &updateInfo)

	// Extract id of the entry from the http request
	vars := mux.Vars(r)
	id := vars["id"]

	// Nothing personal
	if id == "4" {
		fmt.Fprintf(w, "You cannot change Galatasaray, please do not attempt")
		return
	}

	for index, article := range Entries {
		if article.Id == id {
			Entries[index] = updateInfo
		}
	}
}

func deleteEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Perfection
	if id == "4" {
		fmt.Fprintf(w, "You cannot change Galatasaray, please do not attempt")
		return
	}

	for index, article := range Entries {
		if article.Id == id {
			// Love doing this
			Entries = append(Entries[:index], Entries[index+1:]...)
		}
	}

}

func handleRequests() {
	// Creating a mux router which handles different Http Request verbs
	// StrictSlash TRUE is routing something/path and something/path/ to same place.
	dislikedRouter := mux.NewRouter().StrictSlash(true)

	// CRUD := Create, Read, Update, Delete
	// 1st end-point
	dislikedRouter.HandleFunc("/", homePage)

	// POST request (Crud)
	dislikedRouter.HandleFunc("/entry", createNewEntry).Methods("POST")

	// GET requests (cRud)
	dislikedRouter.HandleFunc("/all", returnAllEntries)
	dislikedRouter.HandleFunc("/entry/{id}", returnSingleEntry)

	// UPDATE request (crUd)
	dislikedRouter.HandleFunc("/entry/{id}", updateEntry).Methods("PUT")

	// DELETE request (cruD)
	dislikedRouter.HandleFunc("/entry/{id}", deleteEntry).Methods("DELETE")

	// If connection is down somehow, error will be displayed and os.Exit(1) will occur
	log.Fatal(http.ListenAndServe(":10000", dislikedRouter))
}

func main() {
	fmt.Println("Hey, your localhost is working. Cool!")

	//Four sample randomInfoAboutMe is created
	Entries = []randomInfoAboutMe{
		{Id: "1", Title: "From", Content: "Eskisehir", Desc: "Born and raised in it, had some rough winters but wolf never forgets."},
		{Id: "2", Title: "Age", Content: "FeelingOld", Desc: "Nobody calls me young anymore, sadge."},
		{Id: "3", Title: "VideoGames", Content: "Divinity Original Sin 2", Desc: "I'm amazed by this master piece."},
		{Id: "4", Title: "Football", Content: "Galatasaray", Desc: "GERCEKLERI TARIH YAZAR TARIHI DE GALATASARAY."},
	}

	//Handling request till connection breaks down
	handleRequests()
}
