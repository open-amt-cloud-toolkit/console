# Step 1: Modules caching
FROM golang:1.22-alpine3.19@sha256:c46c4609d3cc74a149347161fc277e11516f523fd8aa6347c9631527da0b7a56 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Tests
FROM golang:1.22-alpine3.19@sha256:c46c4609d3cc74a149347161fc277e11516f523fd8aa6347c9631527da0b7a56
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64

CMD ["go", "test", "-v", "./integration-test/..."]