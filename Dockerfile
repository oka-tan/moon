FROM golang:1.19-alpine AS build
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build .

FROM alpine:3.17.3
WORKDIR /app

RUN apk add tzdata
RUN adduser --disabled-password --no-create-home moon

COPY --from=build /app/moon .

USER moon
CMD ./moon
