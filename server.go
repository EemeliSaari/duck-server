package main


import(
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux" // go get github.com/gorilla/mux
	"github.com/rs/cors" // go get github.com/rs/cors
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

						// Go's rutine to construct the json
type Sighting struct{
	ID string 			`json:"id, omitempty"`
	Specie string		`json:"species, omitempty"`
	Description string	`json:"description, omitempty"`
	Date string			`json:"dateTime, omitempty"`
	Count int			`json:"count, omitempty"`
}


type Specie struct{
	Name string			`json:"name, omitempty"`
}


var sightings []Sighting
var species []Specie


func loadData(){
	// Check if the data exists, runs python script otherwise
	files := []string{"species.json", "sightings.json"}
	for _, element := range files{
		if _, err := os.Stat(element); os.IsNotExist(err){
			argList := []string{"http://github.com/vincit/summer-2018.git", "8080", "sightings", "species"}
			cmd := exec.Command("load_data.py", argList...).Run()
			if(cmd != nil){
				fmt.Println("Successfully loaded the data.")
			}
			break
		}
	}
	raw, err := ioutil.ReadFile("sightings.json")
	if err == nil{
		json.Unmarshal(raw, &sightings)
	}
	raw, err = ioutil.ReadFile("species.json")
	if err == nil{
		json.Unmarshal(raw, &species)
	}
}


func getSightings(res http.ResponseWriter, req *http.Request){
	
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(sightings)
}


func getSpecies(res http.ResponseWriter, req *http.Request){

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(species)
}


func postSightings(res http.ResponseWriter, req *http.Request){

	var sighting Sighting
	_ = json.NewDecoder(req.Body).Decode(&sighting)
	sighting.ID = strconv.Itoa(len(sightings) + 1)
	sightings = append(sightings, sighting)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(req.Body)
}


func main(){

	loadData()

	portNum := 8081
	// Check for cmd args for the user port
	if(len(os.Args) > 1){
		if arg, err := strconv.Atoi(os.Args[1]); err == nil{
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

	log.Fatal(http.ListenAndServe(port, handler)) // Prints error if failed
}
