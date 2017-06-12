.PHONY: test clean bin

bin/assume-role: *.go
	go build -o $@ .

bin: bin/assume-role-Linux bin/assume-role-Darwin bin/assume-role-Windows.exe 

bin/assume-role-Linux: *.go
	env GOOS=linux go build -o $@ .
bin/assume-role-Darwin: *.go
	env GOOS=darwin go build -o $@ .
bin/assume-role-Windows.exe: *.go
	env GOOS=windows go build -o $@ .

clean:
	rm -rf bin/*

test:
	go test -race $(shell go list ./... | grep -v /vendor/)
