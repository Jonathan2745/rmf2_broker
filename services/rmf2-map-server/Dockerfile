FROM golang:1.22-bullseye AS builder
WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . ./

# Build static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/rmf2-map-server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/rmf2-map-server /usr/local/bin/rmf2-map-server
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/rmf2-map-server"]
