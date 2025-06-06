package model

import "time"

type Client struct {
	ID        uint `gorm:"primaryKey"`
	Nome      string
	Email     string
	Phone     string
	CreatedAt time.Time
	CPF       string
}

type Product struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Category string
	Price    float64
	Stock    int
}

type Order struct {
	ID         uint `gorm:"primaryKey"`
	ClientID   uint
	OrderDate  time.Time
	Status     string
	TotalValue float64
	Itens      []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	ProductID uint
	Quantity  int
	Product   Product `gorm:"foreignKey:ProductID"`
}

type Payment struct {
	ID          uint `gorm:"primaryKey"`
	OrderID     uint
	Type        string
	Status      string
	PaymentDate time.Time
}
