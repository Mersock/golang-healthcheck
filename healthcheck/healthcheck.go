package healthcheck

import (
	"fmt"
	"net/http"
	"sync"
)

type healthCheck struct {
	Links []string
}

type HealthCheck interface {
	RunHealthCheck()
	CheckLink(link string, wg *sync.WaitGroup)
}

func NewHealthCheck(links []string) HealthCheck {
	return &healthCheck{
		Links: links,
	}
}

func (hc *healthCheck) RunHealthCheck() {
	var wg sync.WaitGroup
	fmt.Println("Perform website checking...")
	for _, link := range hc.Links {
		wg.Add(1)
		go hc.CheckLink(link, &wg)
	}
	wg.Wait()
	fmt.Println("Done!")

}

func (hc *healthCheck) CheckLink(link string, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down!")
		//c <- link
		return
	}

	fmt.Println(link, "is up!")
	//c <- link
}
