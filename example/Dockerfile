#!/bin/bash

FROM golang:alpine as builder
RUN mkdir /build
WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o api *.go

FROM gcr.io/distroless/static

COPY --from=builder /build/api /

ENV PORT=8080

EXPOSE 8080

CMD ["./api"]