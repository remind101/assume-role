.PHONY: all test clean bins deps docker

bin: bin/assume-role

bins: bin/assume-role-Linux bin/assume-role-Darwin bin/assume-role-Windows.exe 

bin/assume-role: deps *.go
	go build -o $@ .
bin/assume-role-Linux: deps *.go
	env GOOS=linux go build -o $@ .
bin/assume-role-Darwin: deps *.go
	env GOOS=darwin go build -o $@ .
bin/assume-role-Windows.exe: deps *.go
	env GOOS=windows go build -o $@ .

deps:
	go get -v -d ./...

update-deps:
	go get -v -u -d ./...

clean:
	rm -rf bin/*

test: deps
	go test -race ./...

docker:
	docker build --tag assume-role .

docker-test: docker
	docker run -it assume-role /bin/bash -c 'cd src; make test'

