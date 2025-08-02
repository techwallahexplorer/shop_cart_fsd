package handler

import (
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

// Handler is the main entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	router := setupRouter()
	router.ServeHTTP(w, r)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// User endpoints
	router.POST("/api/users", createNewUser)
	router.GET("/api/users", listAllUsers)
	router.POST("/api/users/login", handleUserLogin)

	// Item endpoints
	router.POST("/api/items", createNewItem)
	router.GET("/api/items", listAllItems)

	// Cart endpoints (protected)
	cartGroup := router.Group("/api/carts")
	cartGroup.Use(AuthMiddleware())
	{
		cartGroup.POST("", addItemToCart)
		cartGroup.GET("", fetchCartItems)
	}

	// Order endpoints (protected)
	orderGroup := router.Group("/api/orders")
	orderGroup.Use(AuthMiddleware())
	{
		orderGroup.POST("", createOrder)
		orderGroup.GET("", orderHistoryList)
	}

	return router
}
