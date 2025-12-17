# Go Testify
A rewrite of Testify in Golang. Work in Progress. Hack for Hack Week 2025.

## Task
This project supports executing commands with `task`.

## Build
```
go build .
```

## Test
```
go test ./...
```

## Report
This produces a JUnit-style report.

You need to have [go-junit-report](https://github.com/jstemmer/go-junit-report) installed:  `go install github.com/jstemmer/go-junit-report/v2@latest`
```
go test -v ./... 2>&1 | go-junit-report.exe -set-exit-code > report.xml
```
