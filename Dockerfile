FROM golang:1.14

FROM scratch
EXPOSE 8080

# Copy CA certificates to prevent x509: certificate signed by unknown authority errors
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY s3www /s3www

ENTRYPOINT ["/s3www"]
