FROM golang:1.24-alpine AS builder

WORKDIR /usr/app

COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/app/go-app .

FROM alpine:latest
COPY --from=builder /usr/app/go-app /usr/app/go-app

EXPOSE 8080
ENTRYPOINT ["/usr/app/go-app"]