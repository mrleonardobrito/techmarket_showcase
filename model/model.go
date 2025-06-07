package model

import "time"

type Client struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	Nome      string    `gorm:"column:nome"`
	Email     string    `gorm:"column:email"`
	Phone     string    `gorm:"column:telefone"`
	CreatedAt time.Time `gorm:"column:data_cadastro"`
	CPF       string    `gorm:"column:cpf"`
}

type Product struct {
	ID       uint    `gorm:"primaryKey;column:id"`
	Name     string  `gorm:"column:nome"`
	Category string  `gorm:"column:categoria"`
	Price    float64 `gorm:"column:preco"`
	Stock    int     `gorm:"column:estoque"`
}

type Order struct {
	ID         uint        `gorm:"primaryKey;column:id"`
	ClientID   uint        `gorm:"column:id_cliente"`
	OrderDate  time.Time   `gorm:"column:data_pedido"`
	Status     string      `gorm:"column:status"`
	TotalValue float64     `gorm:"column:valor_total"`
	Itens      []OrderItem `gorm:"-"`
}

type OrderItem struct {
	OrderID   uint    `gorm:"primaryKey;column:id_pedido"`
	ProductID uint    `gorm:"primaryKey;column:id_produto"`
	Quantity  int     `gorm:"column:quantidade"`
	Product   Product `gorm:"-"`
}

type Payment struct {
	ID          uint      `gorm:"primaryKey;column:id"`
	OrderID     uint      `gorm:"column:id_pedido"`
	Type        string    `gorm:"column:tipo"`
	Status      string    `gorm:"column:status"`
	PaymentDate time.Time `gorm:"column:data_pagamento"`
}
