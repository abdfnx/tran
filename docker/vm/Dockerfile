FROM alpine:latest

RUN apk update && apk upgrade && apk add --no-cache ca-certificates

COPY tran /usr/bin/tran

ENTRYPOINT ["/usr/bin/tran"]
