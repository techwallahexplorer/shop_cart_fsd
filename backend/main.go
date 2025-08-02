package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// In-memory database using maps
var (
	users     = make(map[uint]*User)
	items     = make(map[uint]*Item)
	carts     = make(map[uint]*Cart)
	cartItems = make(map[uint]*CartItem)
	orders    = make(map[uint]*Order)
	orderItems = make(map[uint]*OrderItem)
	usersByUsername = make(map[string]*User)
	usersByToken = make(map[string]*User)
	nextUserID uint = 1
	nextItemID uint = 1
	nextCartID uint = 1
	nextCartItemID uint = 1
	nextOrderID uint = 1
	nextOrderItemID uint = 1
	dbMutex sync.RWMutex
)

func main() {
	log.Println("Starting shopping cart backend with in-memory database...")

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// User endpoints
	router.POST("/users", createNewUser)
	router.GET("/users", listAllUsers)
	router.POST("/users/login", handleUserLogin)

	// Item endpoints
	router.POST("/items", createNewItem)
	router.GET("/items", listAllItems)

	// Cart endpoints (protected)
	cartGroup := router.Group("/carts")
	cartGroup.Use(AuthMiddleware())
	{
		cartGroup.POST("", addItemToCart)
		cartGroup.GET("", fetchCartItems)
	}

	// Order endpoints (protected)
	orderGroup := router.Group("/orders")
	orderGroup.Use(AuthMiddleware())
	{
		orderGroup.POST("", createOrder)
		orderGroup.GET("", orderHistoryList)
	}

	log.Println("Server starting on :8080")
	router.Run(":8080")
}
