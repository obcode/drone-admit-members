FROM golang:alpine as builder
RUN apk add -U --no-cache ca-certificates
RUN go build

FROM alpine:3.6
EXPOSE 3000

ENV GODEBUG netdns=go

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /drone-admit-members /bin/

ENTRYPOINT ["/bin/drone-admit-members"]
