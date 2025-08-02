package main

import (
	"time"
)

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
