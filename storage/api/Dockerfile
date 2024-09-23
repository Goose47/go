FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY .. .
RUN go build -o app

FROM alpine

WORKDIR /app
COPY --from=builder /app/app /build/app

CMD ["./app"]