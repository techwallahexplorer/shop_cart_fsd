package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Models
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"`
	Token     string    `gorm:"unique" json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

type Item struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Cart struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"unique;not null" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	CartItems []CartItem `gorm:"foreignKey:CartID" json:"cartItems"`
}

type CartItem struct {
	ID      uint `gorm:"primaryKey" json:"id"`
	CartID  uint `gorm:"not null" json:"cartId"`
	ItemID  uint `gorm:"not null" json:"itemId"`
	Item    Item `gorm:"foreignKey:ItemID" json:"item"`
}

type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"userId"`
	CartID    uint      `gorm:"not null" json:"cartId"`
	CreatedAt time.Time `json:"createdAt"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"orderItems"`
}

type OrderItem struct {
	ID      uint `gorm:"primaryKey" json:"id"`
	OrderID uint `gorm:"not null" json:"orderId"`
	ItemID  uint `gorm:"not null" json:"itemId"`
	Item    Item `gorm:"foreignKey:ItemID" json:"item"`
}

// Request types
type UserRegistrationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
}

type AddItemToCartRequest struct {
	ItemID uint `json:"itemId" binding:"required"`
}

type CreateOrderRequest struct {
	CartID uint `json:"cartId" binding:"required"`
}

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

// Auth Middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}
		dbMutex.RLock()
		user, exists := usersByToken[token]
		dbMutex.RUnlock()
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Attach user info to context for downstream handlers
		c.Set("user", user)
		c.Next()
	}
}

// Utility functions
func hashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func generateToken(userID uint, username string) string {
	h := sha256.New()
	h.Write([]byte(username))
	h.Write([]byte(time.Now().String()))
	h.Write([]byte(string(rune(userID))))
	return hex.EncodeToString(h.Sum(nil))
}

// User handlers
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

// Item handlers
func createNewItem(c *gin.Context) {
	var req ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	item := &Item{
		ID:          nextItemID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CreatedAt:   time.Now(),
	}
	items[nextItemID] = item
	nextItemID++

	c.JSON(http.StatusCreated, item)
}

func listAllItems(c *gin.Context) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	itemList := make([]*Item, 0, len(items))
	for _, item := range items {
		itemList = append(itemList, item)
	}
	c.JSON(http.StatusOK, itemList)
}

// Cart handlers
func addItemToCart(c *gin.Context) {
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userObj.(*User)

	var req AddItemToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Find or create cart for user
	var cart *Cart
	for _, c := range carts {
		if c.UserID == user.ID {
			cart = c
			break
		}
	}
	if cart == nil {
		cart = &Cart{
			ID:        nextCartID,
			UserID:    user.ID,
			CreatedAt: time.Now(),
		}
		carts[nextCartID] = cart
		nextCartID++
	}

	// Add item to cart
	cartItem := &CartItem{
		ID:     nextCartItemID,
		CartID: cart.ID,
		ItemID: req.ItemID,
	}
	cartItems[nextCartItemID] = cartItem
	nextCartItemID++

	c.JSON(http.StatusOK, gin.H{"cartId": cart.ID, "itemId": req.ItemID})
}

func fetchCartItems(c *gin.Context) {
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userObj.(*User)

	dbMutex.RLock()
	defer dbMutex.RUnlock()

	// Find cart for user
	var cart *Cart
	for _, c := range carts {
		if c.UserID == user.ID {
			cart = c
			break
		}
	}

	// If no cart exists, return empty cart response
	if cart == nil {
		c.JSON(http.StatusOK, gin.H{
			"cartId": 0,
			"items":  []interface{}{},
		})
		return
	}

	// Get cart items
	var userCartItems []gin.H
	for _, ci := range cartItems {
		if ci.CartID == cart.ID {
			userCartItems = append(userCartItems, gin.H{
				"itemId": ci.ItemID,
				"id":     ci.ID,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"cartId": cart.ID,
		"items":  userCartItems,
	})
}

// Order handlers
func createOrder(c *gin.Context) {
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userObj.(*User)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Find cart for user
	var cart *Cart
	for _, c := range carts {
		if c.ID == req.CartID && c.UserID == user.ID {
			cart = c
			break
		}
	}
	if cart == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// Get cart items
	var userCartItems []*CartItem
	for _, ci := range cartItems {
		if ci.CartID == cart.ID {
			userCartItems = append(userCartItems, ci)
		}
	}
	if len(userCartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Create order
	order := &Order{
		ID:        nextOrderID,
		UserID:    user.ID,
		CartID:    cart.ID,
		CreatedAt: time.Now(),
	}
	orders[nextOrderID] = order
	nextOrderID++

	// Create order items
	for _, ci := range userCartItems {
		orderItem := &OrderItem{
			ID:      nextOrderItemID,
			OrderID: order.ID,
			ItemID:  ci.ItemID,
		}
		orderItems[nextOrderItemID] = orderItem
		nextOrderItemID++
	}

	c.JSON(http.StatusCreated, order)
}

func orderHistoryList(c *gin.Context) {
	userObj, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userObj.(*User)

	dbMutex.RLock()
	defer dbMutex.RUnlock()

	var userOrders []*Order
	for _, order := range orders {
		if order.UserID == user.ID {
			userOrders = append(userOrders, order)
		}
	}

	c.JSON(http.StatusOK, userOrders)
}
