package handler

import (
	"fmt"
	"net/http"
)

type handler struct {
	ClientID    string
	RedirectUri string
	ClineSecret string
}

type Handler interface {
	RedirectLogin(writer http.ResponseWriter, request *http.Request)
	CallBack(writer http.ResponseWriter, request *http.Request)
}

func NewHandler(clientID string, redirectUri string, clineSecret string) Handler {
	return &handler{
		ClientID:    clientID,
		RedirectUri: redirectUri,
		ClineSecret: clineSecret,
	}
}

func (h *handler) RedirectLogin(writer http.ResponseWriter, request *http.Request) {
	responseType := "code"
	state := "12345abcde"
	scope := "profile%20openid"
	nonce := "09876xyz"
	url := fmt.Sprintf("https://access.line.me/oauth2/v2.1/authorize?response_type=%s&client_id=%s&redirect_uri=%s&state=%s&scope=%s&nonce=%s", responseType, h.ClientID, h.RedirectUri, state, scope, nonce)
	http.Redirect(writer, request, url, http.StatusFound)
}

func (h *handler) CallBack(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("GET params were:", request.URL.Query())
	code := request.URL.Query().Get("code")
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintln(writer, code)
}
