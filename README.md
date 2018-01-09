# duck-server


Eemeli Saari 2018

Duck-server is a simple backend server software implemented with Go programming language. Program was a preliminary assignment for [Vincit](http://www.koodarijahti.fi/) to replace their [server](https://github.com/Vincit/summer-2018/). 

Server responds to http request in a JSON format `Content-Type: application/json`.


## Requirements

Requires [Go](https://golang.org/dl/) and [Git](https://git-scm.com/) for installing. 

Tested with the Go version 1.8.3 windows/amd64


## Install & Run

Inside current %GOPATH%:
```
$ go get -u github.com/EemeliSaari/duck-server
$ cd src/github.com/EemeliSaari/duck-server
$ go run server.go
```


Start with different port (default 8081):
```
$ go run server.go -PORT=<port>
```


## Tests

Tests are done using Go's httptest package and cover all the requests.

Run the tests and see test coverage:
```
$ go test -v
$ go test -cover 
```


## Interface

Initial data in the server is from the old [server](https://github.com/Vincit/summer-2018/). Data has been aquired by using **load_data.py** script.


## Usage

The server process the following requests:

- GET /sightings
    - Returns all the listed sightings
- GET /species
    - Returns all accepted duck species that are supported
- POST /sightings
    - Adds a new sighting to the servers sightings list
    - Generates new ID for the new sighting

