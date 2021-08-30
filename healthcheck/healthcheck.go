package healthcheck

type healthCheck struct {
	Links []string
}

type HealthCheck interface {
}

func NewHealthCheck(links []string) HealthCheck {
	return &healthCheck{
		Links: links,
	}
}
