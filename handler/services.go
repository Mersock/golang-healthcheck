package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type services struct {
	Client HttpClient
}

type Services interface {
	GetToken(code string, payload PayloadLineAuth) (result OauthToken, err error)
	SendReport(accessToken string, payload PayloadSendReport) (statusCode int, err error)
	ResultLogger(statusCode int, payload PayloadSendReport) (message string)
}

func NewServices(client HttpClient) Services {
	return &services{
		Client: client,
	}
}

func (h *services) GetToken(code string, payload PayloadLineAuth) (result OauthToken, err error) {
	endpoint := "https://api.line.me/oauth2/v2.1/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", payload.RedirectUri)
	data.Set("client_id", payload.ClientID)
	data.Set("client_secret", payload.ClientSecret)
	data.Set("code", code)

	r, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return result, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := h.Client.Do(r)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

func (h *services) SendReport(accessToken string, payload PayloadSendReport) (statusCode int, err error) {
	endpoint := "https://backend-challenge.line-apps.com/healthcheck/report"
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return statusCode, err
	}

	r, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return statusCode, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := h.Client.Do(r)
	if err != nil {
		return statusCode, err

	}
	defer res.Body.Close()

	return res.StatusCode, nil
}

func (h *services) ResultLogger(statusCode int, payload PayloadSendReport) (message string) {
	if statusCode == 200 {
		totalTime := payload.TotalTime / int64(time.Millisecond)
		message = "The report healthcheck has been submitted successfully."
		fmt.Println("Checked websites: ", payload.TotalWebsites)
		fmt.Println("Successful websites: ", payload.SuccessLists)
		fmt.Println("Failure websites:: ", payload.FailureLists)
		fmt.Println("Total times to finished checking websites:", totalTime, "ms")
	} else {
		message = "Failed to submit healthcheck report."
		fmt.Println(message)
	}

	return message
}
