# --- STAGE 1: BUILDER ---
FROM golang:1.24-alpine AS builder

WORKDIR /usr/app

COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/app/go-app .


# --- STAGE 2: FINAL (YANG SUDAH DIPERBAIKI) ---
FROM alpine:latest

# 1. Tetapkan direktori kerja di image final (praktik terbaik)
WORKDIR /usr/app

# 2. Salin binary dari builder stage ke direktori kerja saat ini (.)
COPY --from=builder /usr/app/go-app .

# 3. TAMBAHKAN INI: Salin .env dari builder stage ke direktori kerja saat ini (.)
COPY --from=builder /usr/app/.env .


EXPOSE 3000

# 4. Karena WORKDIR sudah diatur, kita bisa memanggil binary secara relatif
ENTRYPOINT ["./go-app"]