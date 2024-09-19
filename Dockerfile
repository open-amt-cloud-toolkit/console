# Step 1: Modules caching
FROM golang:1.23-alpine3.20@sha256:ac67716dd016429be8d4c2c53a248d7bcdf06d34127d3dc451bda6aa5a87bc06 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN apk add --no-cache git
RUN go mod download

# Step 2: Builder
FROM golang:1.23-alpine3.20@sha256:ac67716dd016429be8d4c2c53a248d7bcdf06d34127d3dc451bda6aa5a87bc06 as builder
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