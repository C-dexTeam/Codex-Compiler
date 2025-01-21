# 1. Aşama: Geliştirme ortamı oluşturma
FROM golang:1.23-alpine3.19 AS builder

# Gerekli araçları yükleme
RUN mkdir /app
WORKDIR /app
ENV CGO_ENABLED=1

# Gerekli paketleri yükle (Ruby, Rust, GCC, vs.)
RUN apk update && apk add --no-cache gcc musl-dev git ruby ruby-dev rust cargo

# Air kurulumunu yap
RUN go install github.com/air-verse/air@latest

# Proje bağımlılıklarını yükle
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodunu kopyala
COPY . .

# Air konfigürasyon dosyasını kopyala (varsa)
COPY air.toml .

# 2. Aşama: Çalıştırma
FROM golang:1.23-alpine3.19

# Gerekli diller ve araçlar (Ruby, Rust) tekrar ekleniyor
RUN apk update && apk add --no-cache gcc musl-dev git ruby ruby-dev rust cargo

# Uygulama klasörü oluştur
RUN mkdir /app
WORKDIR /app

# Builder aşamasından gerekli dosyaları al
COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /app /app

# Çalıştırma komutu (Air ile hot-reload)
ENTRYPOINT ["air", "-c", "air.toml"]
