.PHONY: test clean bin

bin/assume-role: *.go
	go build -o $@ .

bin: bin/assume-role-Linux bin/assume-role-Darwin bin/assume-role-Windows.exe 

bin/assume-role-Linux: *.go
	docker run --env GOOS=linux --volume ${PWD}:/go/assume-role --workdir /go/assume-role golang:1.13 go build -mod=vendor -o ./bin/assume-role-Linux
bin/assume-role-Darwin: *.go
	docker run --env GOOS=darwin --volume ${PWD}:/go/assume-role --workdir /go/assume-role golang:1.13 go build -mod=vendor -o ./bin/assume-role-Darwin
bin/assume-role-Windows.exe: *.go
	docker run --env GOOS=windows --volume ${PWD}:/go/assume-role --workdir /go/assume-role  golang:1.13 go build -mod=vendor -o ./bin/assume-role-Windows.exe

clean:
	rm -rf bin/*

test:
	docker run --env TZ=America/New_York --volume ${PWD}:/go/assume-role --workdir /go/assume-role  golang:1.13 go test -mod=vendor -v .
