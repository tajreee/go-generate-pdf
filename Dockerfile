# --- STAGE 1: BUILDER ---
FROM golang:1.24-alpine AS builder

WORKDIR /usr/app

COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/app/go-app .


# --- STAGE 2: FINAL (YANG SUDAH DIPERBAIKI) ---
FROM alpine:latest

RUN apk update && apk add --no-cache \
    chromium \
    ttf-freefont \
    udev \
    xvfb

# Tetapkan direktori kerja di image final (praktik terbaik)
WORKDIR /usr/app

# Salin binary dari builder stage ke direktori kerja saat ini (.)
COPY --from=builder /usr/app/go-app .

# TAMBAHKAN INI: Salin .env dari builder stage ke direktori kerja saat ini (.)
COPY --from=builder /usr/app/.env .


EXPOSE 3000

# Praktik terbaik: Jalankan sebagai user non-root
# '-D' berarti tidak membuat password
RUN adduser -D appuser
USER appuser

# Karena WORKDIR sudah diatur, kita bisa memanggil binary secara relatif
ENTRYPOINT ["./go-app"]