package stream

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/wilsonth122/money-tracker-api/pkg/auth"
	"github.com/wilsonth122/money-tracker-api/pkg/model"
)

type authClient struct {
	userID string
	conn   *websocket.Conn
}

type tokenString struct {
	Token string `json:"token"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var clients = make(map[authClient]bool)
var broadcast = make(chan *model.Expense)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Init() {
	go expense_stream()
}

func Writer(expense *model.Expense) {
	broadcast <- expense
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %s", err)
	}

	err = ws.WriteMessage(websocket.TextMessage, []byte("Waiting for AuthToken"))
	if err != nil {
		log.Printf("Websocket write error: %s", err)
		ws.Close()
		return
	}

	// Wait for client to send users auth token before continuing
	token := tokenString{}
	err = ws.ReadJSON(&token)
	if err != nil {
		log.Printf("Websocket read error: %s", err)
		ws.Close()
		return
	}

	tk, err := auth.ParseToken(token.Token)
	if err != nil {
		log.Printf("Websocket auth error: %s", err)
		ws.Close()
		return
	}

	var client authClient
	client.userID = tk.UserID
	client.conn = ws

	// Register client
	clients[client] = true

	err = ws.WriteMessage(websocket.TextMessage, []byte("Connected"))
	if err != nil {
		log.Printf("Websocket write error: %s", err)
		ws.Close()
	}

	log.Println("Connected")

	clientDone := make(chan struct{})
	go ping(client, clientDone)
	read(client, clientDone)
}

func expense_stream() {
	for {
		expense := <-broadcast

		// Send updated expense to clients that have authenticated as the current user
		for client := range clients {
			if expense.UserID == client.userID {
				err := client.conn.WriteJSON(expense)
				if err != nil {
					log.Printf("Stream, Websocket error: %s", err)
					client.conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func ping(client authClient, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := client.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Printf("Ping, Websocket error: %s", err)
				client.conn.Close()
				delete(clients, client)
				return
			}
		case <-done:
			return
		}
	}
}

func read(client authClient, done chan struct{}) {
	defer client.conn.Close()
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Read, Websocket error: %s", err)
			client.conn.Close()
			delete(clients, client)
			close(done)
			return
		}
	}
}
