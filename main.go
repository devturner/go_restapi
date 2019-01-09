package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

type Application struct {
	ID          string `validate:"nonzero"`
	Title       string `validate:"nonzero"`
	Version     string `validate:"nonzero"`
	Maintainers []struct {
		Name  string `validate:"nonzero"`
		Email string `validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	}
	Company     string `validate:"nonzero"`
	Website     string `validate:"nonzero"`
	Source      string `validate:"nonzero"`
	License     string `validate:"nonzero"`
	Description string `validate:"nonzero"`
}

var appl []Application

// todo: Had to ad ID to YAML file, but that might not be needed if I
// can come up with a better way to search what was provided.
// Need to think about it.

func GetApplicationMetadataEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range appl {
		if item.ID == params["id"] {
			yaml.NewEncoder(w).Encode(item)
			return
		}
	}
	yaml.NewEncoder(w).Encode(&Application{})
}

func GetApplicationsMetadataEndpoint(w http.ResponseWriter, req *http.Request) {
	yaml.NewEncoder(w).Encode(appl)
}

// todo: error handling
func CreateApplicationMetadataEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var newAppl Application
	_ = yaml.NewDecoder(req.Body).Decode(&newAppl)
	newAppl.ID = params["id"]

	if errs := validator.Validate(newAppl); errs != nil {
		http.Error(w, "Field error", http.StatusForbidden)
		return
	}
	appl = append(appl, newAppl)
	yaml.NewEncoder(w).Encode(appl)
}

func DeleteApplicationMetadataEndpoint(w http.ResponseWriter, req *http.Request) {

}

func main() {
	router := mux.NewRouter()

	var appl1 Application
	reader, _ := os.Open("workingExample.yaml")
	buf, _ := ioutil.ReadAll(reader)
	yaml.Unmarshal(buf, &appl1)
	if errs := validator.Validate(appl1); errs != nil {
		fmt.Printf("Field error")
		return
	}

	appl = append(appl, appl1)

	router.HandleFunc("/applications", GetApplicationsMetadataEndpoint).Methods("GET")
	router.HandleFunc("/applications/{id}", GetApplicationMetadataEndpoint).Methods("GET")
	router.HandleFunc("/applications/{id}", CreateApplicationMetadataEndpoint).Methods("POST")
	router.HandleFunc("/applications", GetApplicationsMetadataEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
