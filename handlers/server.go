package handlers

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"freecommunication/models"
)

func CreateServer(w http.ResponseWriter, r *http.Request) {
	var server models.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	db := r.Context().Value("db").(*gorm.DB)
	userID := r.Context().Value("userId").(uint)
	server.UserID = userID

	if err := db.Create(&server).Error; err != nil {
		http.Error(w, "Failed to create server", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(server)
}

func GetServers(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	userID := r.Context().Value("userId").(uint)

	var servers []models.Server
	db.Where("user_id = ?", userID).Find(&servers)
	json.NewEncoder(w).Encode(servers)
}