package main

import (
	"freecommunication/handlers"
	"freecommunication/models"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	db, err := gorm.Open(sqlite.Open("discord.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&models.User{}, &models.Server{}, &models.Channel{}, &models.Message{})

	r := mux.NewRouter()
	r.HandleFunc("/api/auth/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", handlers.Login).Methods("POST")
	r.HandleFunc("/api/servers", handlers.CreateServer).Methods("POST")
	r.HandleFunc("/api/servers", handlers.GetServers).Methods("GET")
	r.HandleFunc("/api/servers/{id}/channels", handlers.CreateChannel).Methods("POST")
	r.HandleFunc("/api/servers/{id}/channels", handlers.GetChannels).Methods("GET")
	r.HandleFunc("/api/channels/{id}/messages", handlers.GetMessages).Methods("GET")
	r.HandleFunc("/api/ws/{channelId}", handlers.HandleWebSocket)

	var port = ":8000"
	log.Println("Server running on " + port)
	http.ListenAndServe(port, r)
}