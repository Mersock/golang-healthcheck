package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type handler struct {
	ClientID    string
	RedirectUri string
	ClineSecret string
}

type Handler interface {
	RedirectLogin(writer http.ResponseWriter, request *http.Request)
	CallBack(writer http.ResponseWriter, request *http.Request)
	getToken(code string) (result OauthToken)
}

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
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
	code := request.URL.Query().Get("code")
	token := h.getToken(code)

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintln(writer, token.AccessToken)
}

func (h *handler) getToken(code string) (result OauthToken) {
	endpoint := "https://api.line.me/oauth2/v2.1/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", h.RedirectUri)
	data.Set("client_id", h.ClientID)
	data.Set("client_secret", h.ClineSecret)
	data.Set("code", code)

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	return result
}
