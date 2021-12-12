//Go HTTP router based on a table of regexes
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}
type ctxKey struct{}

var routes = []route{
	newRoute("GET", "/", homePage),
	newRoute("GET", "/all", returnAllEntries),
	newRoute("POST", "/entry", createNewEntry),

	// Below with input requests, I couldnt make it work
	//newRoute("GET", "/entry/{id}", returnSingleEntry),
	//newRoute("PUT", "/entry/{id}", updateEntry),
	//newRoute("DELETE", "/entry/{id}", deleteEntry),
}

// getField is useful for parsing the request into parameters and pass them to functions
// func getField(r *http.Request, index int) string {
// 	fields := r.Context().Value(ctxKey{}).([]string)
// 	return fields[index]
// }

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

// newRoute is the connection between functions and expressions in the request
func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

// Serve is just pure magic, it just works.
func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
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

	http.HandleFunc("/", Serve)

	http.ListenAndServe(":10000", nil)
}
