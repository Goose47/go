FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY .. .
RUN go build -o app main.go

FROM alpine

WORKDIR /app
COPY --from=builder /app/main /build/app

CMD ["./hello"]