# Dockerfile untuk Realtime ERC20 Transfer Listener (Golang)

# Gunakan image base golang
FROM golang:1.21

# Set working directory
WORKDIR /app

# Copy semua file ke container
COPY . .

# Unduh dependensi
RUN go mod tidy

# Build binary
RUN go build -o erc20-listener .

# Expose port untuk WebSocket
EXPOSE 8080

# Jalankan binary
CMD ["./erc20-listener"]
