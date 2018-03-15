FROM golang:latest

ARG SRV_HOST="0.0.0.0:8000"
ARG SRV_SUBDOMAIN=""

ARG DB_HOST="localhost"
ARG DB_NAME="testdatabase"
ARG DB_USER=""
ARG DB_PWD=""

ARG AUTH_REDIRECT_URL="http://localhost:8000/oauth2callback"
ARG AUTH_CLIENTSECRET="cartracker_api9527"
ARG AUTH_CLIENTID="cartracker-api"
ARG AUTH_COOKIESECRET="cartracker_cookie"

ENV SRV_HOST=$SRV_HOST
ENV SRV_SUBDOMAIN=$SRV_SUBDOMAIN

ENV DB_HOST=$DB_HOST
ENV DB_NAME=$DB_NAME
ENV DB_USER=$DB_USER
ENV DB_PWD=$DB_PWD

ENV AUTH_REDIRECT_URL=$AUTH_REDIRECT_URL
ENV AUTH_CLIENTSECRET=$AUTH_CLIENTSECRET
ENV AUTH_CLIENTID=$AUTH_CLIENTID
ENV AUTH_COOKIESECRET=$AUTH_COOKIESECRET

#disable crosscompiling 
#ENV CGO_ENABLED=0

#compile linux only
#ENV GOOS=linux

COPY . /usr/local/go/src/cartracker.api
WORKDIR /usr/local/go/src/cartracker.api 

#VOLUME /app/settings

ADD VERSION .

RUN go get ./...

#RUN cp settings/*.json /app/settings

#RUN go build -ldflags '-w -s' -a -installsuffix cgo -o main
RUN go build -o main . 

EXPOSE 8000

ENTRYPOINT ["./main"]