package main

import (
	"fmt"
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
	Port            = "8080"
	AppURl          = "http://localhost:" + Port
	LineClientID    = "1656368826"
	LineRedirectUri = AppURl + "/callback"
	LineSecret      = "1ebd163f0ff50083b449d00715603630"
	Timeout         = 10
)

func init() {
	os.Setenv("TZ", "Asia/Bangkok")
}

func main() {
	filename := os.Args[1]
	var wg sync.WaitGroup
	var mutex sync.Mutex
	router := mux.NewRouter()
	client := &http.Client{
		Timeout: Timeout * time.Second,
	}
	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%v", Port),
		WriteTimeout: Timeout * time.Second,
		ReadTimeout:  Timeout * time.Second,
	}

	reader := readcsv.NewReadCSV(filename)
	links, err := reader.ReaderCSV()
	if err != nil {
		log.Fatalf("ReaderCSV Error: %v", err)
	}

	baseUrl := fmt.Sprintf("%v", AppURl)
	hc := healthcheck.NewHealthCheck(links, &wg, client, &mutex, baseUrl)
	sendReport := hc.RunHealthCheck()

	lineAuth := handler.PayloadLineAuth{
		ClientID:     LineClientID,
		RedirectUri:  LineRedirectUri,
		ClientSecret: LineSecret,
	}
	h := handler.NewHandler(lineAuth, sendReport, client)
	router.HandleFunc("/", h.RedirectLogin).Methods(http.MethodGet)
	router.HandleFunc("/callback", h.CallBack).Methods(http.MethodGet)

	log.Fatal(server.ListenAndServe())
}
