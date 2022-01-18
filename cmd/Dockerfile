#Build
FROM golang:1.17.3-alpine3.15 AS build

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./. .
RUN go build ./cmd/main.go

#Environment
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/main .

CMD ["./main", "./wait"]