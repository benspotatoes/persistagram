FROM golang:1.11.5

RUN mkdir -p $GOPATH/src/github.com/benspotatoes/persistagram

COPY . $GOPATH/src/github.com/benspotatoes/persistagram/

ENV DB_ACCESS_TOKEN=dropbox

WORKDIR $GOPATH/src/github.com/benspotatoes/persistagram
