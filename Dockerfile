FROM alpine:3.16.2 as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY legitify /legitify
ENTRYPOINT ["/legitify"]
