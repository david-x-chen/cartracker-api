FROM golang:latest

RUN mkdir -p /app/settings

ADD . /usr/local/go/src/cartracker.api
WORKDIR /usr/local/go/src/cartracker.api 

ADD VERSION .

RUN go get ./...

RUN go build -o main . 

RUN cp settings/*.json /app/settings

ENTRYPOINT ["./main"]