FROM golang:1.12-stretch as build

WORKDIR /src
COPY . .

RUN make test bin
RUN make bins

FROM debian:stretch

COPY --from=build /src/bin/* /usr/local/bin/
RUN chmod 555 /usr/local/bin/assume-role*
