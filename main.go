package main

import (
	"github.com/gorilla/mux"
	"golang-healthcheck/handler"
	"net/http"
)

const (
	LineClientID    = "1656368826"
	LineRedirectUri = "http://localhost:8080/callback"
	LineSecret      = "1ebd163f0ff50083b449d00715603630"
)

func main() {
	router := mux.NewRouter()
	lineAuth := handler.PayloadLineAuth{
		ClientID:     LineClientID,
		RedirectUri:  LineRedirectUri,
		ClientSecret: LineSecret,
	}
	sendReport := handler.PayloadSendReport{
		TotalWebsites: 1000000,
		SuccessLists:  2,
		FailureLists:  1,
		TotalTime:     5000,
	}
	h := handler.NewHandler(lineAuth, sendReport)

	router.HandleFunc("/", h.RedirectLogin).Methods(http.MethodGet)
	router.HandleFunc("/callback", h.CallBack).Methods(http.MethodGet)

	http.ListenAndServe(":8080", router)
}
