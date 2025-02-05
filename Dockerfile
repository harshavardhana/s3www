# Copyright 2021 Harshavardhana
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.22.0 AS builder

WORKDIR /app

RUN \
    mkdir -p /app/licenses && \
    curl -s -q https://raw.githubusercontent.com/harshavardhana/s3www/master/CREDITS -o /app/licenses/CREDITS && \
    curl -s -q https://raw.githubusercontent.com/harshavardhana/s3www/master/LICENSE -o /app/licenses/LICENSE

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o s3www .

FROM scratch

EXPOSE 8080

# Copy CA certificates to prevent x509: certificate signed by unknown authority errors
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/licenses /licenses

COPY --from=builder /app/s3www /s3www

ENTRYPOINT ["/s3www"]
