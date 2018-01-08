package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gorilla/mux" // go get github.com/gorilla/mux
	"github.com/rs/cors"     // go get github.com/rs/cors
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

func loadData() {
	// Check if the data exists, runs python script otherwise
	files := []string{"rsc/species.json", "rsc/sightings.json"}
	for _, element := range files {
		_, err := os.Stat(element)

		if err != nil {

			// Ask user if they want to fetch the original server data
			fmt.Println("Do you want to run python script\nto load the data?\n(Y/N)")
			var out string
			fmt.Scanln(&out)

			if out == "Y" {
				argList := []string{"load_data.py", "http://github.com/vincit/summer-2018.git", "8080", "sightings", "species"}
				cmd := exec.Command("python", argList...)
				err := cmd.Start()
				if err != nil {
					fmt.Println("Failed to start the script")
					return
				}
				err = cmd.Wait()
				if err != nil {
					fmt.Println("Failed to run the script")
				}
				os.RemoveAll("summer-2018") // remove the temporary folder
				fmt.Println("Script completed")
			} else {
				fmt.Println("Server will start without data.")
				break
			}
		}
	}
	raw, err := ioutil.ReadFile("rsc/sightings.json")
	if err == nil {
		json.Unmarshal(raw, &sightings)
	}
	raw, err = ioutil.ReadFile("rsc/species.json")
	if err == nil {
		json.Unmarshal(raw, &species)
	}
}

func getSightings(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(sightings)
}

func getSpecies(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(species)
}

func postSightings(res http.ResponseWriter, req *http.Request) {

	var sighting Sighting
	_ = json.NewDecoder(req.Body).Decode(&sighting)
	sighting.ID = strconv.Itoa(len(sightings) + 1)
	sightings = append(sightings, sighting)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(req.Body)
}

func main() {

	loadData()

	portNum := 8081
	// Check for cmd args for the user port
	if len(os.Args) > 1 {
		if arg, err := strconv.Atoi(os.Args[1]); err == nil {
			portNum = arg
		}
	}
	port := fmt.Sprint(":", portNum)

	// create HTTP request multiplexer
	router := mux.NewRouter()
	router.HandleFunc("/sightings", getSightings).Methods("GET")
	router.HandleFunc("/species", getSpecies).Methods("GET")
	router.HandleFunc("/sightings", postSightings).Methods("POST")

	// toggle CORS options
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Access-Control-Allow-Origin", "Origin, X-Requested-With, Content-Type, Accept", "*"},
	})
	handler := c.Handler(router)

	fmt.Println("Server listening port", portNum)

	log.Fatal(http.ListenAndServe(port, handler)) // Start the actual server
}
