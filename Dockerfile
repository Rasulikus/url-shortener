FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/url-shortener

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /out/app /app/app

EXPOSE 8081
ENTRYPOINT ["/app/app"]