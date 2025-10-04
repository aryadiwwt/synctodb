# Tahap 1: Build Environment (Builder)
# Menggunakan image Go resmi yang berbasis Alpine Linux yang ringan.
FROM golang:1.21-alpine AS builder

# Set direktori kerja di dalam container
WORKDIR /app

# Salin file dependensi terlebih dahulu untuk memanfaatkan Docker cache layer
COPY go.mod go.sum ./

# Unduh semua dependensi
RUN go mod download

# Salin sisa source code aplikasi
COPY . .

# Kompilasi aplikasi.
# CGO_ENABLED=0 membuat binary yang tidak bergantung pada library C dari host.
# Ini penting untuk membuat image final yang super kecil.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# ---

# Tahap 2: Production Environment (Final Image)
# Menggunakan image Alpine Linux yang sangat kecil sebagai dasar.
FROM alpine:latest

# Set direktori kerja
WORKDIR /app

# Salin binary yang sudah dicompile dari tahap 'builder'
COPY --from=builder /app/main .

# Perintah yang akan dijalankan saat container启动
# Menggunakan format exec agar menjadi proses utama (PID 1).
ENTRYPOINT ["/app/main"]