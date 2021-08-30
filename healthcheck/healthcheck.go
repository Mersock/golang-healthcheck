package healthcheck

import (
	"fmt"
	"golang-healthcheck/handler"
	"net/http"
	"sync"
	"time"
)

type healthCheck struct {
	Links     []string
	WaitGroup *sync.WaitGroup
	Client    *http.Client
}

type CountStatus struct {
	Success int
	Failure int
}

type HealthCheck interface {
	RunHealthCheck() handler.PayloadSendReport
	checkLink(link string, success, failure *int)
}

func NewHealthCheck(links []string, waitGroup *sync.WaitGroup, client *http.Client) HealthCheck {
	return &healthCheck{
		Links:     links,
		WaitGroup: waitGroup,
		Client:    client,
	}
}

func (hc *healthCheck) RunHealthCheck() handler.PayloadSendReport {
	fmt.Println("Perform website checking...")
	var success, failure int
	start := time.Now()
	for _, link := range hc.Links {
		hc.WaitGroup.Add(1)
		go hc.checkLink(link, &success, &failure)
	}
	hc.WaitGroup.Wait()
	totalTime := int(time.Since(start).Seconds())
	fmt.Println("Done!")

	return handler.PayloadSendReport{
		TotalWebsites: len(hc.Links),
		SuccessLists:  success,
		FailureLists:  failure,
		TotalTime:     totalTime,
	}
}

func (hc *healthCheck) checkLink(link string, success, failure *int) {
	defer hc.WaitGroup.Done()
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	_, err := client.Get(link)
	if err != nil {
		*failure++
		return
	}
	*success++
}
