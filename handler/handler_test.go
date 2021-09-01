package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	LineClientID    = "1656368826"
	LineRedirectUri = "http://localhost:8080/callback"
	LineSecret      = "1ebd163f0ff50083b449d00715603630"
)

var client = &http.Client{}

var lineAuth = PayloadLineAuth{
	ClientID:     LineClientID,
	RedirectUri:  LineRedirectUri,
	ClientSecret: LineSecret,
}

var payloadReport = PayloadSendReport{
	TotalWebsites: 8,
	SuccessLists:  6,
	FailureLists:  2,
	TotalTime:     5000000000,
}

type ClientMockOK struct{}
type ClientMockFail struct{}

func (c *ClientMockOK) Do(req *http.Request) (*http.Response, error) {
	mockResponse := `{"status":"ok"}`
	t := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockResponse)
	}))
	defer t.Close()
	return http.Get(t.URL)
}

func (c *ClientMockFail) Do(req *http.Request) (*http.Response, error) {
	mockResponse := `{"status":"ok"}`
	t := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, mockResponse)
	}))
	defer t.Close()
	return http.Get(t.URL)
}

func TestRedirectLogin(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	res := w.Result()
	defer res.Body.Close()

	h := NewHandler(lineAuth, payloadReport, client)
	h.RedirectLogin(w, req)
	t.Logf("reditect %+v", res)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, Actual %v", http.StatusOK, res.StatusCode)
	}
}

func TestCallbackSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/callback", nil)
	w := httptest.NewRecorder()
	res := w.Result()
	defer res.Body.Close()

	clientMock := &ClientMockOK{}
	h := NewHandler(lineAuth, payloadReport, clientMock)
	h.CallBack(w, req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, Actual %v", http.StatusOK, res.StatusCode)
	}
}

func TestCallbackFail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/callback", nil)
	w := httptest.NewRecorder()
	res := w.Result()
	defer res.Body.Close()

	clientMock := &ClientMockFail{}
	h := NewHandler(lineAuth, payloadReport, clientMock)
	h.CallBack(w, req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %v, Actual %v", http.StatusOK, res.StatusCode)
	}
}
