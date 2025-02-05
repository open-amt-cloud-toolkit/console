# Step 1: Modules caching
FROM golang:1.23-alpine3.20@sha256:22caeb4deced0138cb4ae154db260b22d1b2ef893dde7f84415b619beae90901 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN apk add --no-cache git
RUN go mod download

# Step 2: Builder
FROM golang:1.23-alpine3.20@sha256:22caeb4deced0138cb4ae154db260b22d1b2ef893dde7f84415b619beae90901 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o /bin/app ./cmd/app
RUN mkdir -p /.config/device-management-toolkit
# Step 3: Final
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/internal/app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /.config/device-management-toolkit /.config/device-management-toolkit
CMD ["/app"]