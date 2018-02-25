FROM golang:1.9-alpine
FROM golang:1.9.3-alpine3.7

ADD . /go/src/github.com/banzaicloud/IngressAuthGenerator
WORKDIR /go/src/github.com/banzaicloud/IngressAuthGenerator
RUN go build -o /bin/IngressAuthGenerator .

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=0 /bin/IngressAuthGenerator /bin
ENTRYPOINT ["/bin/IngressAuthGenerator"]