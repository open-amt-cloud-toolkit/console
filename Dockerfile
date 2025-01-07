# Step 1: Modules caching
FROM golang:1.23-alpine3.20@sha256:6a8532e5441593becc88664617107ed567cb6862cb8b2d87eb33b7ee750f653c AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN apk add --no-cache git
RUN go mod download

# Step 2: Builder
FROM golang:1.23-alpine3.20@sha256:6a8532e5441593becc88664617107ed567cb6862cb8b2d87eb33b7ee750f653c AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/internal/app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app"]