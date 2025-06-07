package seed

import (
	"math/rand"
	"techmarket_showcase/model"
	"time"
)

var orderStatus = []string{
	"Em Processamento",
	"Aguardando Pagamento",
	"Pago",
	"Em Separação",
	"Em Transporte",
	"Entregue",
	"Cancelado",
}

func generateOrderItems(productCount int, maxItems int) []model.OrderItem {
	numItems := rand.Intn(maxItems) + 1
	items := make([]model.OrderItem, numItems)

	usedProducts := make(map[uint]bool)

	for i := 0; i < numItems; i++ {
		var productID uint
		for {
			productID = uint(rand.Intn(productCount) + 1)
			if !usedProducts[productID] {
				usedProducts[productID] = true
				break
			}
		}

		items[i] = model.OrderItem{
			ProductID: productID,
			Quantity:  rand.Intn(5) + 1,
		}
	}

	return items
}

func GenerateOrders(count int, clientCount int, productCount int) []model.Order {
	orders := make([]model.Order, count)
	maxItemsPerOrder := 10

	for i := 0; i < count; i++ {
		orderDate := time.Now().Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour)

		items := generateOrderItems(productCount, maxItemsPerOrder)

		var totalValue float64
		for _, item := range items {
			productPrice := 100.0 + rand.Float64()*900.0
			totalValue += productPrice * float64(item.Quantity)
		}

		orders[i] = model.Order{
			ClientID:   uint(rand.Intn(clientCount) + 1),
			OrderDate:  orderDate,
			Status:     orderStatus[rand.Intn(len(orderStatus))],
			TotalValue: float64(int(totalValue*100)) / 100,
			Itens:      items,
		}
	}

	return orders
}
