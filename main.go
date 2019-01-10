package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	appl    []Application
	results []Application
)

// validEmail validate email address
func validEmail(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validator.ErrUnsupported
	}
	re, err := regexp.Compile("[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$")
	if err != nil {
		return errors.New("Enter a valid email address")
	}
	if !re.MatchString(st.String()) {
		return errors.New("Enter a valid email address")
	}
	return nil
}

// Application MetaData record type
type Application struct {
	ID          string `validate:"nonzero"`
	Title       string `validate:"nonzero"`
	Version     string `validate:"nonzero"`
	Maintainers []struct {
		Name  string `validate:"nonzero"`
		Email string `validate:"nonzero, validEmail"`
	} `validate:"nonzero"`
	Company     string `validate:"nonzero"`
	Website     string `validate:"nonzero"`
	Source      string `validate:"nonzero"`
	License     string `validate:"nonzero"`
	Description string `validate:"nonzero"`
}

// GetApplicationMetadataEndpoint Search by ID for record
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

// SearchApplicationMetadataEndpoint Search application titles
func SearchApplicationMetadataEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range appl {
		if strings.Contains(item.Title, params["key"]) {
			results = append(results, item)
		} else {
			io.WriteString(w, "No results for that key")
		}
	}
	yaml.NewEncoder(w).Encode(results)
	results = results[:0]
}

// GetApplicationsMetadataEndpoint Return all entries
func GetApplicationsMetadataEndpoint(w http.ResponseWriter, req *http.Request) {
	yaml.NewEncoder(w).Encode(appl)
}

// CreateApplicationMetadataEndpoint Create a new application record
func CreateApplicationMetadataEndpoint(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)
	var newAppl Application
	_ = yaml.NewDecoder(req.Body).Decode(&newAppl)
	newAppl.ID = params["id"]

	err := validator.Validate(newAppl)
	if err == nil {
		appl = append(appl, newAppl)
		yaml.NewEncoder(w).Encode(appl)
	} else {
		errs := err.(validator.ErrorMap)
		var errOuts []string
		for f, e := range errs {
			errOuts = append(errOuts, fmt.Sprintf("\t - %s (%v)\n", f, e))
		}
		sort.Strings(errOuts)
		http.Error(w, "Invalid due to field(s):", http.StatusForbidden)
		for _, str := range errOuts {
			io.WriteString(w, str)
		}
	}
}

// DeleteApplicationMetadataEndpoint Delete a record
func DeleteApplicationMetadataEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range appl {
		if item.ID == params["id"] {
			appl = append(appl[:index], appl[index+1:]...)
			break
		}
	}
	yaml.NewEncoder(w).Encode(appl)
}

func main() {
	// register custom validation
	validator.SetValidationFunc("validEmail", validEmail)

	router := mux.NewRouter()

	// test with workingExample.yaml & brokenExample.yaml
	var appl1 Application
	// reader, _ := os.Open("brokenExample1.yaml")
	reader, _ := os.Open("workingExample1.yaml")
	buf, _ := ioutil.ReadAll(reader)
	yaml.Unmarshal(buf, &appl1)
	if errs := validator.Validate(appl1); errs != nil {
		fmt.Printf("Field error")
		return
	}
	// if valid, add to prject at startup
	appl = append(appl, appl1)

	// routes for GET all, GET one by ID, POST, GET Search Title, DELETE one by ID
	router.HandleFunc("/applications", GetApplicationsMetadataEndpoint).Methods("GET")
	router.HandleFunc("/applications/{id}", GetApplicationMetadataEndpoint).Methods("GET")
	router.HandleFunc("/search/{key}", SearchApplicationMetadataEndpoint).Methods("GET")
	router.HandleFunc("/new/{id}", CreateApplicationMetadataEndpoint).Methods("POST")
	router.HandleFunc("/delete/{id}", DeleteApplicationMetadataEndpoint).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
