package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"freecommunication/models"
)

func CreateChannel(w http.ResponseWriter, r *http.Request) {
	var channel models.Channel
	if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	serverID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}
	channel.ServerID = uint(serverID)

	db := r.Context().Value("db").(*gorm.DB)
	if err := db.Create(&channel).Error; err != nil {
		http.Error(w, "Failed to create channel", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(channel)
}

func GetChannels(w http.ResponseWriter, r *http.Request) {
	serverID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}
	db := r.Context().Value("db").(*gorm.DB)

	var channels []models.Channel
	db.Where("server_id = ?", serverID).Find(&channels)
	json.NewEncoder(w).Encode(channels)
}