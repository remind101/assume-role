FROM golang:1.12-stretch

WORKDIR /app
COPY . ./src

RUN cd src && make bin
RUN mkdir bin && cp -ai src/bin/* bin/ && chmod +x bin/*
