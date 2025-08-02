package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
}

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
