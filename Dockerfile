FROM alpine:3.5

RUN apk --no-cache add ca-certificates

ADD storage /bin/server