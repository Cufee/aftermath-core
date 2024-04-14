# Download localization files
FROM node:21-alpine as assets

WORKDIR /app

COPY accent.json aftermath-core.json ./

ARG ACCENT_API_KEY

RUN mkdir -p ./internal/core/localization/resources
RUN npm install -g accent-cli && accent export

# Build the app binary
FROM golang:1.22.1-alpine as builder

WORKDIR /app 

COPY . .
COPY --from=assets /app/internal/core/localization/resources ./internal/core/localization/resources

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o binary .

# Make a scratch container with required files and binary
FROM scratch

WORKDIR /app

ENV TZ=Europe/Berlin
ENV ZONEINFO=/zoneinfo.zip
COPY --from=builder /app/binary /usr/bin/
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["binary"]