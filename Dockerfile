FROM golang:1.26.2-alpine3.23 AS build

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

WORKDIR /go/src
COPY ./go.* .golangci.yaml /go/src
RUN tree /go/src && \
    go mod download -x && \
    mkdir -p /go/src/bin && \
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

COPY . /go/src
RUN tree /go/src

RUN go build -o bin/ees-link cmd/ees-link/main.go

FROM build AS test

WORKDIR /go/src
RUN mkdir -p /go/src/bin/tests && \
    go test ./... -o /dev/null -c && \
    sh -c 'go test ./... -json -count=1 -timeout 30s -cover -coverprofile=bin/tests/cover.txt || true' > bin/tests/tests.json

FROM alpine:3.23.0 AS final

COPY --from=build ./go/src/bin/ees-link ./

RUN chmod +x /ees-link

WORKDIR /
CMD ["./ees-link"]