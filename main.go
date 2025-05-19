package main

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow any origin
	},
}

type TransferEvent struct {
	From  common.Address `json:"from"`
	To    common.Address `json:"to"`
	Value *big.Int       `json:"value"`
}

var wsClients = make(map[*websocket.Conn]bool)
var addClient = make(chan *websocket.Conn)
var removeClient = make(chan *websocket.Conn)
var broadcast = make(chan TransferEvent)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to Ethereum node
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/" + os.Getenv("INFURA_PROJECT_ID"))
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	// Load ERC20 ABI
	erc20ABI, err := abi.JSON(strings.NewReader(`[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	go handleClients()
	go listenTransferEvents(client, erc20ABI)

	http.HandleFunc("/ws", handleWebSocket)
	log.Println("WebSocket server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func listenTransferEvents(client *ethclient.Client, contractABI abi.ABI) {
	query := ethereum.FilterQuery{
		Topics: [][]common.Hash{{contractABI.Events["Transfer"].ID}},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Println("Subscription error:", err)
		case vLog := <-logs:
			log.Printf("Received log: %v", vLog)
			var event TransferEvent
			err := contractABI.UnpackIntoInterface(&event, "Transfer", vLog.Data)
			if err != nil {
				log.Println("Failed to unpack event:", err)
				continue
			}
			event.From = common.HexToAddress(vLog.Topics[1].Hex())
			event.To = common.HexToAddress(vLog.Topics[2].Hex())
			broadcast <- event
		}
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	addClient <- conn

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			removeClient <- conn
			conn.Close()
			break
		}
	}
}

func handleClients() {
	for {
		select {
		case conn := <-addClient:
			wsClients[conn] = true
		case conn := <-removeClient:
			delete(wsClients, conn)
		case event := <-broadcast:
			msg, _ := json.Marshal(event)
			for conn := range wsClients {
				err := conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					removeClient <- conn
					conn.Close()
				}
			}
		}
	}
}
