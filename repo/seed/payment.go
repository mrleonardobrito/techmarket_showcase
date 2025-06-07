package seed

import (
	"math/rand"
	"techmarket_showcase/model"
	"time"
)

var (
	paymentTypes = []string{
		"PIX",
		"Cartão de Crédito",
		"Cartão de Débito",
		"Boleto",
		"Transferência Bancária",
	}

	paymentStatus = []string{
		"Aprovado",
		"Pendente",
		"Recusado",
		"Em Processamento",
		"Estornado",
	}

	statusProbability = map[string]int{
		"Aprovado":         70,
		"Pendente":         10,
		"Recusado":         5,
		"Em Processamento": 10,
		"Estornado":        5,
	}
)

func generatePaymentStatus() string {
	roll := rand.Intn(100)
	accumulated := 0

	for status, prob := range statusProbability {
		accumulated += prob
		if roll < accumulated {
			return status
		}
	}

	return "Aprovado"
}

func GeneratePayments(orderCount int, count int) []model.Payment {
	payments := make([]model.Payment, count)

	for i := 0; i < count; i++ {
		paymentDate := time.Now().Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour)

		payments[i] = model.Payment{
			OrderID:     uint(rand.Intn(orderCount) + 1),
			Type:        paymentTypes[rand.Intn(len(paymentTypes))],
			Status:      generatePaymentStatus(),
			PaymentDate: paymentDate,
		}
	}

	return payments
}
