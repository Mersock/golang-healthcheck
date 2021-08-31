package main

import (
	"github.com/gorilla/mux"
	"golang-healthcheck/handler"
	"golang-healthcheck/healthcheck"
	"golang-healthcheck/readcsv"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	LineClientID    = "1656368826"
	LineRedirectUri = "http://localhost:8080/callback"
	LineSecret      = "1ebd163f0ff50083b449d00715603630"
)

func init() {
	os.Setenv("TZ", "Asia/Bangkok")
}

func main() {
	var wg sync.WaitGroup
	var client = &http.Client{
		Timeout: 10 * time.Second,
	}
	var mutex sync.Mutex
	router := mux.NewRouter()

	reader := readcsv.NewReadCSV("test.csv")
	links, err := reader.ReaderCSV()
	if err != nil {
		log.Fatalf("ReaderCSV Error: %v", err)
	}

	hc := healthcheck.NewHealthCheck(links, &wg, client, &mutex)
	sendReport := hc.RunHealthCheck()

	lineAuth := handler.PayloadLineAuth{
		ClientID:     LineClientID,
		RedirectUri:  LineRedirectUri,
		ClientSecret: LineSecret,
	}

	h := handler.NewHandler(lineAuth, sendReport, client)
	router.HandleFunc("/", h.RedirectLogin).Methods(http.MethodGet)
	router.HandleFunc("/callback", h.CallBack).Methods(http.MethodGet)

	http.ListenAndServe(":8080", router)
}
