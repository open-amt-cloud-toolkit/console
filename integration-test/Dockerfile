# Step 1: Modules caching
FROM golang:1.22-alpine3.19@sha256:f1fe698725f6ed14eb944dc587591f134632ed47fc0732ec27c7642adbe90618 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Tests
FROM golang:1.22-alpine3.19@sha256:f1fe698725f6ed14eb944dc587591f134632ed47fc0732ec27c7642adbe90618
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64

CMD ["go", "test", "-v", "./integration-test/..."]