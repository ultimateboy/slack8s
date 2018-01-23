FROM alpine:3.4
RUN apk add --no-cache ca-certificates
ADD slack8s /slack8s

ENTRYPOINT ["/slack8s"]
