package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"net/http"
	"time"
	"freecommunication/models"
)

type Claims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	db := r.Context().Value("db").(*gorm.DB)
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	token := generateToken(user.ID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	db := r.Context().Value("db").(*gorm.DB)
	var user models.User
	if err := db.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := generateToken(user.ID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func generateToken(userID uint) string {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("secret"))
	return tokenString
}