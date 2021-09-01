# golang-healthcheck

## Package
* [gorilla/mux](https://github.com/gorilla/mux) implements a request router and dispatcher for matching incoming requests to their respective handler.
* Golang standard libraries and including `net/http` package.

## Install
 * Required to install [golang 1.17](https://golang.org/dl/) .
 * Support to run with [docker](https://www.docker.com/products/docker-desktop) container 

## Instructions
1. `Perform website checking...` is step of execute healthcheck list websites from file `test.csv`

2. When get output `Done!` Please open web browser and go to http://localhost:8080/ to login to submit healthcheck report.
3. `Total times to finished checking websits` show with output by milliseconds and send to `https://backend-challenge.line-apps.com/healthcheck/report` by unix nano


## Examples compile the packages and run.

```
$ cd golang-healthcheck
$ go build -o go-healthcheck 
$ ./go-healthcheck test.csv 
Perform website checking...
Done!

Please open web browser and go to http://localhost:8080/ to login to submit healthcheck report.

Checked websites:  8
Successful websites:  6
Failure websites::  2
Total times to finished checking websites: 17644 ms
```

## Examples run by docker container
```
$ cd golang-healthcheck
$ docker build -t go-healthcheck .
$ docker run -p 8080:8080 go-healthcheck test.csv
Perform website checking...
Done!

Please open web browser and go to http://localhost:8080/ to login to submit healthcheck report.

Checked websites:  8
Successful websites:  6
Failure websites::  2
Total times to finished checking websites: 17159 ms

```

## Tests
```
$ go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```