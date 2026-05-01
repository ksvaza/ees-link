# Stage 1: Build Go server
FROM golang:1.26.2-alpine3.23 AS build

WORKDIR /go/src

COPY . /go/src

RUN tree /go/src

RUN go build -o ees-link cmd/ees-link/main.go

# Stage 2: Final stage
FROM alpine:3.23.0 AS final

COPY --from=build ./go/src/ees-link ./

RUN chmod +x /ees-link

WORKDIR /
CMD ["./ees-link"]