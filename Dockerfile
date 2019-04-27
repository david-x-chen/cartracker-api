FROM golang:latest as builder

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

WORKDIR /usr/local/go/src/cartracker.api 

COPY . .

# Download dependencies
RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/cartrackerapi .

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/cartrackerapi .

EXPOSE 8000

ENTRYPOINT ["./cartrackerapi"]