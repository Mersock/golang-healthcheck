package healthcheck

import (
	"fmt"
	"net/http"
	"sync"
)

type healthCheck struct {
	Links     []string
	WaitGroup *sync.WaitGroup
}

type HealthCheck interface {
	RunHealthCheck()
	CheckLink(link string)
}

func NewHealthCheck(links []string, waitGroup *sync.WaitGroup) HealthCheck {
	return &healthCheck{
		Links:     links,
		WaitGroup: waitGroup,
	}
}

func (hc *healthCheck) RunHealthCheck() {
	fmt.Println("Perform website checking...")
	for _, link := range hc.Links {
		hc.WaitGroup.Add(1)
		go hc.CheckLink(link)
	}
	hc.WaitGroup.Wait()
	fmt.Println("Done!")

}

func (hc *healthCheck) CheckLink(link string) {
	defer hc.WaitGroup.Done()
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "might be down!")
		//c <- link
		return
	}

	fmt.Println(link, "is up!")
	//c <- link
}
