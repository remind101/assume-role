bin/assume-role: *.go
	go build -o bin/assume-role .

test:
	go test -race $(shell go list ./... | grep -v /vendor/)
