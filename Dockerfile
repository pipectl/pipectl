FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY pipectl /usr/local/bin/pipectl
ENTRYPOINT ["pipectl"]
