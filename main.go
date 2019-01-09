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

	"github.com/gorilla/mux"
	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

func validEmail(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validate.ErrUnsupported
	}
	if st.String() != regexp.MatchString("^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$") {
		return errors.New("Enter a valid email address")
	}
	return nil
}

validate.SetValidationFunc("validEmail", validEmail)

type Application struct {
	ID          string `validate:"nonzero"`
	Title       string `validate:"nonzero"`
	Version     string `validate:"nonzero"`
	Maintainers []struct {
		Name  string `validate:"nonzero"`
		Email string `validate:"nonzero, validEmail"`
	}
	// "regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Company     string `validate:"nonzero"`
	Website     string `validate:"nonzero"`
	Source      string `validate:"nonzero"`
	License     string `validate:"nonzero"`
	Description string `validate:"nonzero"`
}

var appl []Application

// todo: Had to ad ID to YAML file, but that might not be needed if I
// can come up with a better way to search what was provided.
// This might help:
// https://stackoverflow.com/questions/38654383/how-to-search-for-an-element-in-a-golang-slice
// https://github.com/lithammer/fuzzysearch

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
		http.Error(w, "Invalid due to fields:", http.StatusForbidden)
		for _, str := range errOuts {
			io.WriteString(w, str)
		}
	}
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
