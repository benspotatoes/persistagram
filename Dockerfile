FROM golang:1.11.5

RUN mkdir -p $GOPATH/src/github.com/benspotatoes/persistagram
RUN mkdir -p /opt/persistagram/data

COPY . $GOPATH/src/github.com/benspotatoes/persistagram/

ENV LIKED_FILE=/liked.txt
ENV SAVE_DIRECTORY=/opt/persistagram/data

WORKDIR $GOPATH/src/github.com/benspotatoes/persistagram
