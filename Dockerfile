FROM golang:1.10.1

RUN mkdir -p $GOPATH/src/github.com/benspotatoes/persistagram

COPY . $GOPATH/src/github.com/benspotatoes/persistagram/

ENV DB_ACCESS_TOKEN=dropbox

WORKDIR $GOPATH/src/github.com/benspotatoes/persistagram

RUN go build -o gram main.go
