# Stage 1: Build Go server
FROM --platform=linux/arm64 golang:1.23.4-alpine3.21 AS build

ENV GOOS=linux GOARCH=arm64

COPY ./go.* /go/src
COPY ./cmd /go/src/cmd
#COPY ./master /go/src/master
#COPY ./public /go/src/public

WORKDIR /go/src
RUN go build -o ees-link cmd/ees-link/main.go

# Stage 2: Final stage
FROM --platform=linux/arm/v8 alpine:3.21.0 AS final

COPY --from=build ./go/src/ees-link ./
#COPY --from=build ./go/src/public ./public

RUN chmod +x /ees-link

WORKDIR /
CMD ["./ees-link"]