# Build
FROM golang:1.18-alpine as build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
COPY cmd ./
COPY pb ./
COPY pkg ./
COPY internal ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build ./app -o ./cmd/app

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ./app /app

EXPOSE 5555

ENTRYPOINT ["/app"]