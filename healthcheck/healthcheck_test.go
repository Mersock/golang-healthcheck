package healthcheck

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	Port   = "8080"
	AppURl = "http://localhost:" + Port
)

var wg sync.WaitGroup
var mutex sync.Mutex
var client = &http.Client{
	Timeout: 10 * time.Second,
}

func TestRunHealthCheck(t *testing.T) {
	baseUrl := fmt.Sprintf("%v", AppURl)
	links := []string{"https://line.me/th/", "https://reqres.in/api/users?delay=15", "https://winning.co.th"}
	hc := NewHealthCheck(links, &wg, client, &mutex, baseUrl)
	sendReport := hc.RunHealthCheck()
	//t.Logf("xx %+v", sendReport)
	if sendReport.TotalWebsites != len(links) {
		t.Errorf("Expected %v, Actual %v", len(links), sendReport.TotalWebsites)
	}
	if sendReport.SuccessLists != 1 {
		t.Errorf("Expected %v, Actual %v", 1, sendReport.SuccessLists)
	}
	if sendReport.FailureLists != 2 {
		t.Errorf("Expected %v, Actual %v", 2, sendReport.FailureLists)
	}
	if sendReport.TotalTime == 0 {
		t.Errorf("Expected greater than %v, Actual %v", 0, sendReport.TotalTime)
	}
}
