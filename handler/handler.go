package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type handler struct {
	PayloadLineAuth   PayloadLineAuth
	PayloadSendReport PayloadSendReport
	Client            *http.Client
}

type PayloadLineAuth struct {
	ClientID     string
	RedirectUri  string
	ClientSecret string
}

type PayloadSendReport struct {
	TotalWebsites int `json:"total_websites"`
	SuccessLists  int `json:"success"`
	FailureLists  int `json:"failure"`
	TotalTime     int `json:"total_time"`
}

type Handler interface {
	RedirectLogin(writer http.ResponseWriter, request *http.Request)
	CallBack(writer http.ResponseWriter, request *http.Request)
	getToken(code string) (result OauthToken)
	sendReport(accessToken string) (statusCode int)
}

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

func NewHandler(auth PayloadLineAuth, report PayloadSendReport, client *http.Client) Handler {
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
	token := h.getToken(code)
	statusCode := h.sendReport(token.AccessToken)
	var text string
	if statusCode == 200 {
		text = "The report healthcheck has been submitted successfully."
		fmt.Println("Checked webistes: ", h.PayloadSendReport.TotalWebsites)
		fmt.Println("Successful websites: ", h.PayloadSendReport.SuccessLists)
		fmt.Println("Failure websites:: ", h.PayloadSendReport.FailureLists)
		fmt.Println("Total times to finished checking website:", h.PayloadSendReport.TotalTime, "sec")
	} else {
		text = "Failed to submit healthcheck report."
		fmt.Println(text)
	}

	writer.WriteHeader(statusCode)
	fmt.Fprintln(writer, text)
}

func (h *handler) getToken(code string) (result OauthToken) {
	endpoint := "https://api.line.me/oauth2/v2.1/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", h.PayloadLineAuth.RedirectUri)
	data.Set("client_id", h.PayloadLineAuth.ClientID)
	data.Set("client_secret", h.PayloadLineAuth.ClientSecret)
	data.Set("code", code)

	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := h.Client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	return result
}

func (h *handler) sendReport(accessToken string) (statusCode int) {
	endpoint := "https://backend-challenge.line-apps.com/healthcheck/report"
	jsonData, err := json.Marshal(h.PayloadSendReport)
	if err != nil {
		log.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := h.Client.Do(r)
	if err != nil {
		log.Fatal(err)

	}
	defer res.Body.Close()

	return res.StatusCode
}
