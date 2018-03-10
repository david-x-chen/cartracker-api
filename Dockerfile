FROM golang:latest

ADD . /usr/local/go/src/cartracker.api
WORKDIR /usr/local/go/src/cartracker.api 

VOLUME /app/settings

ADD VERSION .

RUN go get ./...

RUN cp settings/*.json /app/settings

RUN go build -o main . 

ENTRYPOINT ["./main"]