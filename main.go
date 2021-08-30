package main

import (
	"github.com/gorilla/mux"
	"golang-healthcheck/healthcheck"
	"golang-healthcheck/readcsv"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	r := readcsv.NewReadCSV("test.csv")
	links := r.ReaderCSV()
	healthcheck.NewHealthCheck(links)

	http.ListenAndServe(":8080", router)
}
