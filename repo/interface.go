package repo

import (
	"techmarket_showcase/model"
	"time"
)

type TechMarketRepository interface {
	BatchCreateClient(clients []model.Client) error
	BatchCreateProduct(products []model.Product) error
	BatchCreateOrder(orders []model.Order) error
	BatchCreateOrderItem(orderItems []model.OrderItem) error
	BatchCreatePayment(payments []model.Payment) error

	GetClientByEmail(email string) (model.Client, error)
	GetProductByCategory(category string) ([]model.Product, error)
	GetDeliveredOrdersByClient(clientID uint) ([]model.Order, error)
	Get5MostSoldProducts() ([]model.Product, error)
	GetLastMonthPixPayments() ([]model.Payment, error)
	GetClientTotalSpentByPeriod(clientID uint, startDate time.Time, endDate time.Time) (float64, error)
}
