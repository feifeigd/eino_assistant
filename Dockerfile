FROM golang:1.27-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/app ./cmd/einoagent

FROM alpine:3.21
RUN apk --no-cache add ca-certificates redis \
    && update-ca-certificates

WORKDIR /app
COPY .env ./
COPY data data

COPY --from=builder /out/app ./app
EXPOSE 8080
CMD ["./app"]
