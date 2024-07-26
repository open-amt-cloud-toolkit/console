# Step 1: Modules caching
FROM golang:1.22-alpine3.20@sha256:0d3653dd6f35159ec6e3d10263a42372f6f194c3dea0b35235d72aabde86486e as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22-alpine3.20@sha256:0d3653dd6f35159ec6e3d10263a42372f6f194c3dea0b35235d72aabde86486e as builder
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