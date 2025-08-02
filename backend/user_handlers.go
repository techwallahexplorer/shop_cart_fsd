package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserRegistrationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func hashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func createNewUser(c *gin.Context) {
	var req UserRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Check if username already exists
	if _, exists := usersByUsername[req.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	user := &User{
		ID:           nextUserID,
		Username:     req.Username,
		PasswordHash: hashPassword(req.Password),
		CreatedAt:    time.Now(),
	}
	users[nextUserID] = user
	usersByUsername[req.Username] = user
	nextUserID++

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "username": user.Username})
}

func listAllUsers(c *gin.Context) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	userList := make([]*User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	c.JSON(http.StatusOK, userList)
}

func handleUserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	user, exists := usersByUsername[req.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password"})
		return
	}
	if user.PasswordHash != hashPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username/password"})
		return
	}
	token := generateToken(user.ID, user.Username)
	user.Token = token
	usersByToken[token] = user
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func generateToken(userID uint, username string) string {
	h := sha256.New()
	h.Write([]byte(username))
	h.Write([]byte(time.Now().String()))
	h.Write([]byte(string(rune(userID))))
	return hex.EncodeToString(h.Sum(nil))
}
