package handler

import "net/http"

type PayloadSendReport struct {
	TotalWebsites int   `json:"total_websites"`
	SuccessLists  int   `json:"success"`
	FailureLists  int   `json:"failure"`
	TotalTime     int64 `json:"total_time"`
}

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

type PayloadLineAuth struct {
	ClientID     string
	RedirectUri  string
	ClientSecret string
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
