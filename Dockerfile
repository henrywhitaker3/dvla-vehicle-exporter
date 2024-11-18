FROM alpine:3.20.3 AS certs

RUN apk add ca-certificates

FROM scratch
WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY dvla-vehicle-exporter /dvla-vehicle-exporter
USER 65532:65532

ENTRYPOINT ["/dvla-vehicle-exporter"]
