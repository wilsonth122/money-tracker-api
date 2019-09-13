package stream

import (
	"log"
	"net/http"

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
		log.Println(err)
	}

	err = ws.WriteMessage(websocket.TextMessage, []byte("Waiting for AuthToken"))
	if err != nil {
		log.Printf("Websocket error: %s", err)
		ws.Close()
	}

	// Wait for client to send users auth token before continuing
	token := tokenString{}
	err = ws.ReadJSON(&token)
	if err != nil {
		log.Printf("Websocket error: %s", err)
	}

	tk, err := auth.ParseToken(token.Token)
	if err != nil {
		log.Printf("Websocket error: %s", err)
		return
	}

	var client authClient
	client.userID = tk.UserID
	client.conn = ws

	// Register client
	clients[client] = true

	err = ws.WriteMessage(websocket.TextMessage, []byte("Connected"))
	if err != nil {
		log.Printf("Websocket error: %s", err)
		ws.Close()
	}

	log.Println("Connected")
}

func expense_stream() {
	for {
		expense := <-broadcast

		// Send updated expense to clients that have authenticated as the current user
		for client := range clients {
			if expense.UserID == client.userID {
				log.Println("Sending update: " + expense.Title)
				err := client.conn.WriteJSON(expense)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}
