package handlers

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"freecommunication/models"
	"encoding/json"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	channelID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}
	db := r.Context().Value("db").(*gorm.DB)

	var messages []models.Message
	db.Where("channel_id = ?", channelID).Preload("User").Find(&messages)
	json.NewEncoder(w).Encode(messages)
}