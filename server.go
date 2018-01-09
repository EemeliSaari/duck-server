package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" 
	"github.com/rs/cors" 
)

//Sighting struct that contains all the duck sightings
type Sighting struct {
	ID          string `json:"id, omitempty"`
	Specie      string `json:"species, omitempty"`
	Description string `json:"description, omitempty"`
	Date        string `json:"dateTime, omitempty"`
	Count       int    `json:"count, omitempty"`
}

//Specie struct that contains names of all the duck species
type Specie struct {
	Name string `json:"name, omitempty"`
}

var sightings []Sighting
var species []Specie

// LoadData loads the interface data to the arrays
func LoadData() {

	if raw, err := ioutil.ReadFile("rsc/sightings.json"); err == nil{
		json.Unmarshal(raw, &sightings)
	} else{
		fmt.Println("Error reading file: %v", err)
	}
	if raw, err := ioutil.ReadFile("rsc/species.json"); err == nil{
		json.Unmarshal(raw, &species)
	} else{
		fmt.Println("Error reading file: %v", err)
	}
}

// GetSightings returns list of all the sightings as json http response
func GetSightings(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(sightings)
}

// GetSpecies returns list of all the species as json http response
func GetSpecies(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(species)
}

// PostSightings adds sighting to the sightings list
func PostSightings(res http.ResponseWriter, req *http.Request) {

	var sighting Sighting
	_ = json.NewDecoder(req.Body).Decode(&sighting)
	sighting.ID = strconv.Itoa(len(sightings) + 1)
	sightings = append(sightings, sighting)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(req.Body)
}

func MakeHandler() http.Handler{

	// create HTTP request multiplexer
	router := mux.NewRouter()
	router.HandleFunc("/sightings", GetSightings).Methods("GET")
	router.HandleFunc("/species", GetSpecies).Methods("GET")
	router.HandleFunc("/sightings", PostSightings).Methods("POST")

	// toggle CORS options
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Access-Control-Allow-Origin", "Origin, X-Requested-With, Content-Type, Accept", "*"},
	})
	handler := c.Handler(router)

	return(handler)
}

// cmd flag
var Port = flag.String("PORT", "8081", "Listen address")

func main() {

	LoadData()
	flag.Parse()
	address := fmt.Sprint(":", *Port)

	fmt.Println("Server listening port", *Port)

	log.Fatal(http.ListenAndServe(address, MakeHandler())) // Start the actual server
}
