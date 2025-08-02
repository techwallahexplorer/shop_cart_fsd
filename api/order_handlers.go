package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	CartID uint `json:"cartId" binding:"required"`
}

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
