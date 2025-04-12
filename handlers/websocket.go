package handlers

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"freecommunication/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[uint]map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	channelID, err := strconv.Atoi(mux.Vars(r)["channelId"])
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value("userId").(uint)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	if clients[uint(channelID)] == nil {
		clients[uint(channelID)] = make(map[*websocket.Conn]bool)
	}
	clients[uint(channelID)][conn] = true

	go handleMessages(db)

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket read error:", err)
			delete(clients[uint(channelID)], conn)
			break
		}

		msg.ChannelID = uint(channelID)
		msg.UserID = userID
		db.Create(&msg)
		db.Preload("User").First(&msg, msg.ID)
		broadcast <- msg
	}
}

func handleMessages(db *gorm.DB) {
	for msg := range broadcast {
		for conn := range clients[msg.ChannelID] {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println("WebSocket write error:", err)
				conn.Close()
				delete(clients[msg.ChannelID], conn)
			}
		}
	}
}