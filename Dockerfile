FROM golang:1.17-alpine

RUN mkdir "/app"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o go-healthcheck .

EXPOSE 8080

ENTRYPOINT ["./go-healthcheck"]

CMD ["File"]







