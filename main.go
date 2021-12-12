package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	// Using gorilla temporarily due to lack of time
	"github.com/gorilla/mux"
)

// json tags are useful for unmarshalling of request bodies
type randomInfoAboutMe struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Content string `json:"content"`
	Desc    string `json:"desc"`
}

// Entries populated in main, used as a simulation for DB
var Entries []randomInfoAboutMe

// alreadyExist checks if given Id is already in Entries
func alreadyExist(Id string) bool {
	for _, entry := range Entries {
		if entry.Id == Id {
			return true
		}
	}
	return false
}

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

	// If requested id is not found on the Entries
	if !alreadyExist(key) {
		fmt.Fprintf(w, "Requested Id is not found, please request a valid ID")
	}
}

// returnAllEntries returns all entries
func returnAllEntries(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Entries)
}

// createNewEntry handles POST request with a body by appending it to global sturctu Entries
func createNewEntry(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createNewEntry")
	// Get the body of our POST request
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newInfo randomInfoAboutMe

	// JSON to randomInfoAboutMe(My Struct)
	json.Unmarshal(reqBody, &newInfo)

	// Checks request body's Id, if an entry with same ID exists, refuses to create
	// By this way, deleting an Id wont throw an exception due to multiple Id's trying to be deleted at the same time.
	if alreadyExist(newInfo.Id) {
		fmt.Fprintf(w, "You cannot created a new entry that already exist, please use another id.\nYou can GET all and see which id's are already taken.")
		return
	}

	// Update our global Entries array with newly added entry
	Entries = append(Entries, newInfo)

	// This line prints the body of the request it received to inform the client (double-checking)
	json.NewEncoder(w).Encode(newInfo)

	// Sorting Entries after each creation
	sort.Slice(Entries[:], func(i, j int) bool {
		return Entries[i].Id < Entries[j].Id
	})
}

// updateEntry updates only the pointed(by id) struct with the given body information
func updateEntry(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updateEntry")
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

	for index, entry := range Entries {
		if entry.Id == id {
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

	for index, entry := range Entries {
		if entry.Id == id {
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

	// UPDATE request (crUd)
	dislikedRouter.HandleFunc("/entry/{id}", updateEntry).Methods("PUT")

	// DELETE request (cruD)
	dislikedRouter.HandleFunc("/entry/{id}", deleteEntry).Methods("DELETE")

	// GET requests (cRud), careful moving returnSingleEntry to any id requested functions will disable those.
	dislikedRouter.HandleFunc("/all", returnAllEntries)
	dislikedRouter.HandleFunc("/entry/{id}", returnSingleEntry)

	// If connection is down somehow, error will be displayed and os.Exit(1) will occur
	log.Fatal(http.ListenAndServe(":10000", dislikedRouter))
}

func main() {
	fmt.Println("Hey, your localhost is working. Cool!")

	// Sample with 4 elements is created
	Entries = []randomInfoAboutMe{
		{Id: "1", Title: "From", Content: "Eskisehir", Desc: "Born and raised in it, had some rough winters but wolf never forgets."},
		{Id: "2", Title: "Age", Content: "FeelingOld", Desc: "Nobody calls me young anymore, sadge."},
		{Id: "3", Title: "VideoGames", Content: "Divinity Original Sin 2", Desc: "I'm amazed by this master piece."},
		{Id: "4", Title: "Football", Content: "Galatasaray", Desc: "GERCEKLERI TARIH YAZAR TARIHI DE GALATASARAY."},
	}

	// Handling request till connection breaks down
	handleRequests()
}
