package main

import(
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func testGetSightings(t *testing.T){
	
	req, _ := http.NewRequest("GET", "/sightings", nil)
	handler := http.HandlerFunc(GetSightings)

	res := httptest.NewRecorder() 
	handler.ServeHTTP(res, req)

	// Check if the server handled the request 
	if status := res.Code; status != http.StatusOK{
		t.Errorf("Status code %v, excepted %v", status, http.StatusOK)
	}

	var resList []Sighting

	_ = json.NewDecoder(res.Body).Decode(&resList)

	// Check if the response matched the sightings list
	if !reflect.DeepEqual(resList, sightings){
		t.Errorf("Response body isn't valid")
	}
}

func testPostSightings(t *testing.T, id string){

	// Test body, fields are random, since front-end is responsible checking everything but the id
	post := Sighting{ID:id, Specie:"empty", Description:"empty", Date:"empty", Count:42}
	correct := strconv.Itoa(len(sightings) + 1)
	
	var before []Sighting
	copy(before, sightings)

	js, _ := json.Marshal(post)
	body := bytes.NewBuffer(js)

	req, _ := http.NewRequest("POST", "/sightings", body)
	handler := http.HandlerFunc(PostSightings)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	
	// Check if the server handled the request
	if status := res.Code; status != http.StatusOK{
		t.Errorf("Status code %v, excepted %v", status, http.StatusOK)
	}

	// Check if the POST was added to the list
	if len(before) >= len(sightings){
		t.Errorf("Sighting wasn't added to the list")
	}

	correctPost := Sighting{ID:correct, Specie:"empty", Description:"empty", Date:"empty", Count:42}
	added := sightings[len(sightings) - 1]

	// Check if the POST was added correctly
	if added != correctPost{
		t.Errorf("POST wasn't saved correctly")
	}

	// Check if the id was registered correctly
	if added.ID != correct{
		t.Errorf("ID wasn't added correctly")
	}
}

func testGetSpecies(t *testing.T){

	req, _ := http.NewRequest("GET", "/species", nil)
	handler := http.HandlerFunc(GetSpecies)

	res := httptest.NewRecorder() 
	handler.ServeHTTP(res, req)

	// Check if the server handled the request 
	if status := res.Code; status != http.StatusOK{
		t.Errorf("Status code %v, excepted %v", status, http.StatusOK)
	}

	var resList []Specie

	_ = json.NewDecoder(res.Body).Decode(&resList)

	// Check if the response matched the species list
	if !reflect.DeepEqual(resList, species){
		t.Errorf("Response body isn't valid")
	}
}

// TestLoad tests if the load function populated the server lists
func TestLoad(t *testing.T){
	
	LoadData()

	// Simple test if the lists contain anything
	if len(sightings) == 0{
		t.Errorf("Loading sightings failed")
	}
	if len(species) == 0{
		t.Errorf("Loading species failed")
	}
}

// TestServer tests the functionality of the server.
func TestServer(t *testing.T){

	LoadData()
	// Initialize the test server
	ts := httptest.NewServer(MakeHandler())
	url := fmt.Sprint("localhost:", *Port)
	// Change the URL
	ts.URL = url

	defer ts.Close()

	// Run the subtests
	t.Run("GET/sightings", func(t *testing.T){testGetSightings(t)})
	t.Run("POST/sightings/EmptyID", func(t *testing.T){testPostSightings(t, "")})
	t.Run("POST/sightings/NotEmptyID", func(t *testing.T){testPostSightings(t, "3")})
	t.Run("GET/species", func(t *testing.T){testGetSpecies(t)})
}
