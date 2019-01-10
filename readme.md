# Application Metadata API Server
**Author**: Chris Turner

**Version**: 0.1.0 

## Overview

Golang RESTful API server for application metadata compatible with valid YAML data. 

## API Endpoints

    - /applications
    Get list of applications
    - /applications/{id}
    Get one application based on ID search
    - /search/{key}
    Return list of applications by searching titles for keyword
    - /new/{id}
    Persist a new application
    - /delete/{id}
    Delete an application 

 ## Getting Started:
 Clone this repo

 In a terminal instance build & start the server:
    
    - go build && ./restapi

Using Postman, verify the API is working at endpoints:
    
    - GET localhost:8000/applications
    - GET localhost:8000/applications/{id}
    - GET localhost:8000/search/{key}
    - POST localhost:8000/new/{id} 
    - DELETE localhost:8000/delete/{id}

## Dependencies: 
 - gorilla/mux: https://github.com/gorilla/mux
 - validator:  github.com/go-validator/validator/tree/v2
 - YAML: https://github.com/go-yaml/yaml/tree/v2.2.2
