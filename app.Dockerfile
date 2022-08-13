# syntax=docker/dockerfile:1

# Build
FROM golang:1.18-alpine as builder

WORKDIR /builder

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /goapp ./cmd/app/

## Deploy
FROM alpine:3.15

WORKDIR /

COPY --from=builder /goapp /goapp

EXPOSE 8080

ENTRYPOINT ["/goapp"]