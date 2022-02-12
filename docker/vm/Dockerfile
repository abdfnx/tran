FROM alpine:latest

RUN apk update && apk upgrade && apk add --no-cache ca-certificates curl

RUN curl -sL https://cutt.ly/tran-cli | bash

ENTRYPOINT ["/usr/local/bin/tran"]
