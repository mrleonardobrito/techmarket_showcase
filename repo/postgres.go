package repo

import (
	"fmt"
	"techmarket_showcase/config"
	"techmarket_showcase/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ TechMarketRepository = &PostgresRepository{}

type PostgresRepository struct {
	db *gorm.DB
}

func (p *PostgresRepository) BatchCreateClient(clients []model.Client) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		batchSize := 100
		for i := 0; i < len(clients); i += batchSize {
			end := min(i+batchSize, len(clients))

			if err := tx.Table("cliente").Create(clients[i:end]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *PostgresRepository) BatchCreateOrder(orders []model.Order) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		batchSize := 100
		for i := 0; i < len(orders); i += batchSize {
			end := min(i+batchSize, len(orders))

			if err := tx.Table("pedido").Create(orders[i:end]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *PostgresRepository) BatchCreateOrderItem(orderItems []model.OrderItem) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		batchSize := 100
		for i := 0; i < len(orderItems); i += batchSize {
			end := min(i+batchSize, len(orderItems))

			items := make([]map[string]any, len(orderItems[i:end]))
			for j, item := range orderItems[i:end] {
				items[j] = map[string]any{
					"id_pedido":      item.OrderID,
					"id_produto":     item.ProductID,
					"quantidade":     item.Quantity,
					"preco_unitario": item.Product.Price,
				}
			}

			if err := tx.Table("item_pedido").Create(items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *PostgresRepository) BatchCreatePayment(payments []model.Payment) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		batchSize := 100
		for i := 0; i < len(payments); i += batchSize {
			end := min(i+batchSize, len(payments))

			if err := tx.Table("pagamento").Create(payments[i:end]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *PostgresRepository) BatchCreateProduct(products []model.Product) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		batchSize := 100
		for i := 0; i < len(products); i += batchSize {
			end := min(i+batchSize, len(products))

			if err := tx.Table("produto").Create(products[i:end]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (p *PostgresRepository) GetClientByEmail(email string) (model.Client, error) {
	query := `SELECT * FROM cliente WHERE email = ?`
	var client model.Client
	if err := p.db.Raw(query, email).Scan(&client).Error; err != nil {
		return model.Client{}, err
	}
	return client, nil
}

func (p *PostgresRepository) GetProductByCategory(category string) ([]model.Product, error) {
	query := `SELECT * FROM produto WHERE categoria = ?`
	var products []model.Product
	if err := p.db.Raw(query, category).Scan(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (p *PostgresRepository) GetDeliveredOrdersByClient(clientID uint) ([]model.Order, error) {
	query := `SELECT * FROM pedido WHERE id_cliente = ? AND status = 'entregue'`
	var orders []model.Order
	if err := p.db.Raw(query, clientID).Scan(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (p *PostgresRepository) Get5MostSoldProducts() ([]model.Product, error) {
	query := `
		SELECT p.*, COALESCE(SUM(ip.quantidade), 0) as total_vendas
		FROM produto p
		LEFT JOIN item_pedido ip ON p.id = ip.id_produto
		GROUP BY p.id
		ORDER BY total_vendas DESC
		LIMIT 5
	`
	var products []model.Product
	if err := p.db.Raw(query).Scan(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (p *PostgresRepository) GetLastMonthPixPayments() ([]model.Payment, error) {
	query := `SELECT * FROM pagamento WHERE tipo = 'pix' AND data_pagamento >= ? AND data_pagamento <= ?`
	var payments []model.Payment
	if err := p.db.Raw(query, time.Now().AddDate(0, -1, 0), time.Now()).Scan(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (p *PostgresRepository) GetClientTotalSpentByPeriod(clientID uint, startDate time.Time, endDate time.Time) (float64, error) {
	query := `SELECT SUM(valor_total) FROM pagamento WHERE id_cliente = ? AND data_pagamento >= ? AND data_pagamento <= ?`
	var total float64
	if err := p.db.Raw(query, clientID, startDate, endDate).Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func NewPostgresRepository() *PostgresRepository {
	config := config.LoadPostgresConfig()

	db, err := gorm.Open(postgres.Open(config.URI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(fmt.Sprintf("Erro ao conectar com o PostgreSQL: %v", err))
	}

	return &PostgresRepository{db: db}
}
