# 🧠 Realtime ERC20 Transfer Listener (Golang)

Proyek ini adalah backend server **realtime** menggunakan **Golang** yang mendengarkan event `Transfer` dari kontrak ERC-20 di blockchain Ethereum, dan mengirimkan informasi tersebut ke **WebSocket clients** secara langsung. Cocok untuk digunakan sebagai backend **wallet tracker**, **DEX analytics**, atau **notifikasi pengguna**.

---

## 🚀 Fitur Utama

- Terhubung ke Ethereum node melalui WebSocket (Infura atau node lokal)
- Mendengarkan event `Transfer` ERC-20 secara real-time
- Mem-parsing log blockchain dan mengubahnya jadi objek `TransferEvent`
- Menyiarkan event ke semua klien WebSocket yang tersambung
- Skalabel dan ringan, cocok untuk monitoring token, DEX, atau wallet tracker

---

## 🧰 Teknologi yang Digunakan

| Komponen          | Teknologi                         |
|-------------------|-----------------------------------|
| Bahasa            | Go (Golang)                       |
| Ethereum client   | `go-ethereum (ethclient)`         |
| WebSocket server  | `gorilla/websocket`               |
| ABI Parser        | `github.com/ethereum/go-ethereum/accounts/abi` |
| Data format       | JSON                              |

---

## 📦 Instalasi

### 1. Prasyarat
- Go 1.19 atau lebih baru
- Akses ke node Ethereum yang mendukung WebSocket (misalnya [Infura](https://infura.io))

### 2. Clone Proyek
```bash
git clone https://github.com/username/realtime-erc20-listener.git
cd realtime-erc20-listener
```

### 3. Ganti Project ID Infura
Ganti `PROJECT_ID` di file `main.go` dengan project ID Anda di Infura.
```go
client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOUR_INFURA_PROJECT_ID")
```
### 4. Install Dependencies
```bash
go mod download
```
### 5. Run Server
```bash
go run main.go
```
## 🌐 Cara Menggunakan WebSocket
Client dapat terhubung ke server WebSocket dengan URL 
```bash
ws://localhost:8080/ws
```
Setiap kali event Transfer dari token ERC-20 terjadi, client akan menerima data JSON seperti ini:
```json
{
    "From": "0xSenderAddress...",
    "To": "0xReceiverAddress...",
    "Value": "1000000000000000000"
}
```
