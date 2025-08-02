package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AddItemToCartRequest struct {
	ItemID uint `json:"itemId" binding:"required"`
}

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
