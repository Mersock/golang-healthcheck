devup:
	docker compose -f docker-compose.dev.yml  up -d

devdown:
	docker compose -f docker-compose.dev.yml  down

logs:
	docker logs -f golang-healthcheck

tests:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out