FROM golang:1.22.5-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -v -o main main.go

FROM alpine:latest

# TODO: slim this down, chromium is massive
RUN apk --no-cache add ca-certificates chromium font-noto-emoji

WORKDIR /tmp

COPY --from=builder /app/main /app/main

CMD ["/app/main"]