FROM golang:alpine as builder
RUN apk add -U --no-cache ca-certificates
WORKDIR /go/src/github.com/obcode/drone-admit-members
COPY . .
RUN go build

FROM alpine
EXPOSE 3000

ENV GODEBUG netdns=go

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/obcode/drone-admit-members/drone-admit-members /bin/

ENTRYPOINT ["/bin/drone-admit-members"]
