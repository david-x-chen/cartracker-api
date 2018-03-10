FROM golang:latest

ARG mongoHost="192.168.139.218"
ARG serverHost="0.0.0.0:5000"
ARG authRedirectURL="https://dyntechsolution.info/car/oauth2callback"
ARG subLocation="carweb"

ENV mongoHost=${mongoHost}
ENV serverHost=${serverHost}
ENV authRedirectURL=${authRedirectURL}
ENV subLocation=${subLocation}

ADD . /usr/local/go/src/cartracker.api
WORKDIR /usr/local/go/src/cartracker.api 

VOLUME /app/settings

ADD VERSION .

RUN curl http://stedolan.github.io/jq/download/linux64/jq -o /usr/local/bin/jq
RUN chmod a+x /usr/local/bin/jq

# set MongoDB host
RUN cat settings/db_config.json | jq .
RUN jq --arg mhost ${mongoHost} -c '.mongodbhosts = $mhost' settings/db_config.json > tmp.$$.json && mv tmp.$$.json settings/db_config.json
RUN cat settings/db_config.json

# set server host
RUN cat settings/server_config.json | jq .
RUN jq --arg shost ${serverHost} --arg subLoc ${subLocation} -c '.host = $shost | .subLocation = $subLoc' settings/server_config.json > tmp.$$.json && mv tmp.$$.json settings/server_config.json
RUN cat settings/server_config.json

# set auth redirect URL
RUN cat settings/oauth_config.json | jq .
RUN jq --arg authUrl ${authRedirectURL} -c '.redirectUrl = $authUrl' settings/oauth_config.json > tmp.$$.json && mv tmp.$$.json settings/oauth_config.json
RUN cat settings/oauth_config.json

RUN go get ./...

RUN go build -o main . 

RUN cp settings/*.json /app/settings

ENTRYPOINT ["./main"]