#syntax=docker/dockerfile:1
FROM golang:1.26.2-alpine3.23 AS build

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

WORKDIR /go/src
COPY ./go.* .golangci.yaml /go/src
COPY ./.cache /go/pkg/mod
RUN tree /go/src && \
    go env -w GOMODCACHE=/go/pkg/mod

RUN go mod download -x && \
    mkdir -p /go/src/bin && \
    go install github.com/jstemmer/go-junit-report/v2@v2.0.0 && \
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0 && \
    go install github.com/boumenot/gocover-cobertura@v1.4.0

COPY --exclude=.cache . /go/src
RUN tree /go/src

RUN go build -o bin/ees-link cmd/ees-link/main.go

FROM build AS test

WORKDIR /go/src
RUN mkdir -p /go/src/bin/tests && \
    go test ./... -o /dev/null -c && \
    sh -c 'go test ./... -json -count=1 -timeout 30s -cover -coverprofile=bin/tests/coverage.txt || true' > bin/tests/tests.json && \
    ${GOPATH}/bin/go-junit-report -parser gojson < bin/tests/tests.json > bin/tests/tests.xml && \
    ${GOPATH}/bin/gocover-cobertura < bin/tests/coverage.txt > bin/tests/coverage.xml && \
    ${GOPATH}/bin/golangci-lint --config .golangci.yaml run ./... --show-stats=false --output.text.print-issued-lines=false --output.text.colors --issues-exit-code=0 --max-same-issues=0 > bin/tests/linter.txt

FROM alpine:3.23.0 AS final

COPY --from=build ./go/src/bin/ees-link ./

RUN chmod +x /ees-link

WORKDIR /
CMD ["./ees-link"]