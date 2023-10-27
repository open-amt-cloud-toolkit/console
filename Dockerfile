#*********************************************************************
# Copyright (c) Intel Corporation 2023
# SPDX-License-Identifier: Apache-2.0
#*********************************************************************/

#build stage
FROM golang:alpine@sha256:96634e55b363cb93d39f78fb18aa64abc7f96d372c176660d7b8b6118939d97b AS builder
RUN apk add --no-cache git ca-certificates && update-ca-certificates
RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "1000" "scratchuser"
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app -ldflags="-s -w" -v ./cmd/

#final stage
FROM scratch
COPY --from=builder /go/bin/app /app
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER scratchuser
ENTRYPOINT ["/app"]
#LABEL Name=app Version=1.0.0 # Add a label if you wish for your app
LABEL license='SPDX-License-Identifier: Apache-2.0' \
      copyright='Copyright (c) 2023: Intel'
      
EXPOSE 3000
