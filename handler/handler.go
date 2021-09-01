package handler

import (
	"fmt"
	"net/http"
)

type handler struct {
	PayloadLineAuth   PayloadLineAuth
	PayloadSendReport PayloadSendReport
	Client            HttpClient
}

type Handler interface {
	RedirectLogin(writer http.ResponseWriter, request *http.Request)
	CallBack(writer http.ResponseWriter, request *http.Request)
}

func NewHandler(auth PayloadLineAuth, report PayloadSendReport, client HttpClient) Handler {
	return &handler{
		PayloadLineAuth:   auth,
		PayloadSendReport: report,
		Client:            client,
	}
}

func (h *handler) RedirectLogin(writer http.ResponseWriter, request *http.Request) {
	responseType := "code"
	state := "12345abcde"
	scope := "profile%20openid"
	nonce := "09876xyz"
	url := fmt.Sprintf("https://access.line.me/oauth2/v2.1/authorize?response_type=%s&client_id=%s&redirect_uri=%s&state=%s&scope=%s&nonce=%s", responseType, h.PayloadLineAuth.ClientID, h.PayloadLineAuth.RedirectUri, state, scope, nonce)
	http.Redirect(writer, request, url, http.StatusFound)
}

func (h *handler) CallBack(writer http.ResponseWriter, request *http.Request) {
	code := request.URL.Query().Get("code")
	s := NewServices(h.Client)
	token, err := s.GetToken(code, h.PayloadLineAuth)
	if err != nil {
		fmt.Fprintln(writer, "Failed to LINE login ", err)
		return
	}

	statusCode, err := s.SendReport(token.AccessToken, h.PayloadSendReport)
	if err != nil {
		fmt.Fprintln(writer, "Failed to send healthcheck report ", err)
		return
	}

	message := s.ResultLogger(statusCode, h.PayloadSendReport)

	writer.WriteHeader(statusCode)
	fmt.Fprintln(writer, message)
}
