FROM golang:latest as builder

WORKDIR /usr/local/go/src/cartracker.api 

COPY . .

# Download dependencies
RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/cartrackerapi .

######## Start a new stage from scratch #######
FROM alpine:latest  

ENV SRV_HOST="0.0.0.0:8000"
ENV SRV_SUBDOMAIN=""

ENV DB_HOST="localhost"
ENV DB_NAME="testdatabase"
ENV DB_USER=""
ENV DB_PWD=""

ENV AUTH_REDIRECT_URL="http://localhost:8000/oauth2callback"
ENV AUTH_CLIENTSECRET="cartracker_api9527"
ENV AUTH_CLIENTID="cartracker-api"
ENV AUTH_COOKIESECRET="cartracker_cookie"

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/cartrackerapi .

EXPOSE 8000

ENTRYPOINT ["./cartrackerapi"]