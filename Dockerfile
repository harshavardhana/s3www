FROM golang:1.14

ADD go.mod /go/src/github.com/harshavardhana/s3www/go.mod
ADD go.sum /go/src/github.com/harshavardhana/s3www/go.sum
WORKDIR /go/src/github.com/harshavardhana/s3www/
# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download

ADD . /go/src/github.com/harshavardhana/s3www/
WORKDIR /go/src/github.com/harshavardhana/s3www/

ENV CGO_ENABLED=0

RUN go build -ldflags '-w -s' -a -o s3www .

FROM scratch
EXPOSE 8080

COPY --from=0 /go/src/github.com/harshavardhana/s3www/s3www /s3www

ENTRYPOINT ["/s3www"]
