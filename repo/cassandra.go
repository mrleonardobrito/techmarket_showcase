package repo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"techmarket_showcase/config"
	"techmarket_showcase/model"
	"time"

	"github.com/gocql/gocql"
)

var _ TechMarketRepository = &CassandraRepository{}

type CassandraRepository struct {
	db *gocql.Session
}

func (c *CassandraRepository) BatchCreateClient(clients []model.Client) error {
	batchSize := 100
	for i := 0; i < len(clients); i += batchSize {
		end := min(i+batchSize, len(clients))
		batch := c.db.NewBatch(gocql.UnloggedBatch)

		for _, client := range clients[i:end] {
			uuid := gocql.TimeUUID()
			batch.Query(`
				INSERT INTO clientes_por_email (
					email,
					id,
					nome,
					telefone,
					data_cadastro,
					cpf
				) VALUES (?, ?, ?, ?, ?, ?)`,
				client.Email,
				uuid.String(),
				client.Nome,
				client.Phone,
				client.CreatedAt,
				client.CPF,
			)
		}

		if err := c.db.ExecuteBatch(batch); err != nil {
			return err
		}
	}
	return nil
}

func (c *CassandraRepository) BatchCreateProduct(products []model.Product) error {
	batchSize := 100
	for i := 0; i < len(products); i += batchSize {
		end := min(i+batchSize, len(products))
		batch := c.db.NewBatch(gocql.UnloggedBatch)

		for _, product := range products[i:end] {
			uuid := gocql.TimeUUID()

			batch.Query(`
				INSERT INTO produtos_por_categoria (
					categoria,
					preco,
					id_produto,
					nome,
					estoque
				) VALUES (?, ?, ?, ?, ?)`,
				product.Category,
				product.Price,
				uuid.String(),
				product.Name,
				product.Stock,
			)
		}

		if err := c.db.ExecuteBatch(batch); err != nil {
			return err
		}
	}
	return nil
}

func (c *CassandraRepository) BatchCreateOrder(orders []model.Order) error {
	batchSize := 30
	for i := 0; i < len(orders); i += batchSize {
		end := min(i+batchSize, len(orders))
		batch := c.db.NewBatch(gocql.UnloggedBatch)

		for _, order := range orders[i:end] {
			var itensMap []map[string]string
			for _, item := range order.Itens {
				itemMap := map[string]string{
					"product_id": gocql.TimeUUID().String(),
					"quantity":   fmt.Sprintf("%d", item.Quantity),
				}
				itensMap = append(itensMap, itemMap)
			}

			clientID := gocql.TimeUUID()
			orderID := gocql.TimeUUID()

			batch.Query(`
				INSERT INTO pedidos_por_cliente (
					id_cliente,
					data_pedido,
					id_pedido,
					status,
					valor_total,
					itens
				) VALUES (?, ?, ?, ?, ?, ?)`,
				clientID.String(),
				order.OrderDate,
				orderID.String(),
				order.Status,
				order.TotalValue,
				itensMap,
			)
		}

		if err := c.db.ExecuteBatch(batch); err != nil {
			return err
		}
	}
	return nil
}

func (c *CassandraRepository) BatchCreateOrderItem(orderItems []model.OrderItem) error {
	return nil
}

func (c *CassandraRepository) BatchCreatePayment(payments []model.Payment) error {
	batchSize := 100
	for i := 0; i < len(payments); i += batchSize {
		end := min(i+batchSize, len(payments))
		batch := c.db.NewBatch(gocql.UnloggedBatch)

		for _, payment := range payments[i:end] {
			mesAno := payment.PaymentDate.Format("2006-01")
			paymentID := gocql.TimeUUID()
			orderID := gocql.TimeUUID()
			clientID := gocql.TimeUUID()

			batch.Query(`
				INSERT INTO pagamentos_por_tipo_e_mes (
					tipo,
					mes_ano,
					data_pagamento,
					id_pagamento,
					id_pedido,
					id_cliente
				) VALUES (?, ?, ?, ?, ?, ?)`,
				payment.Type,
				mesAno,
				payment.PaymentDate,
				paymentID.String(),
				orderID.String(),
				clientID.String(),
			)
		}

		if err := c.db.ExecuteBatch(batch); err != nil {
			return err
		}
	}
	return nil
}

func (c *CassandraRepository) GetClientByEmail(email string) (model.Client, error) {
	query := `SELECT * FROM clientes_por_email WHERE email = ?`
	var client model.Client
	err := c.db.Query(query, email).Scan(&client)
	if err != nil && err == gocql.ErrNotFound {
		return model.Client{}, nil
	}
	if err != nil {
		return model.Client{}, err
	}
	return client, nil
}

func (c *CassandraRepository) GetProductByCategory(category string) ([]model.Product, error) {
	query := `SELECT * FROM produtos_por_categoria WHERE categoria = ?`
	var products []model.Product
	err := c.db.Query(query, category).Scan(&products)
	if err != nil && err == gocql.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (c *CassandraRepository) GetDeliveredOrdersByClient(clientID uint) ([]model.Order, error) {
	query := `
		SELECT pedido_id, data_pedido, status, valor_total, itens
		FROM pedidos_por_cliente
		WHERE id_cliente = ? AND status = 'entregue'
		ALLOW FILTERING
	`

	var orders []model.Order
	var (
		pedidoIDStr string
		dataPedido  time.Time
		status      string
		valorTotal  float64
		itensJSON   string
	)

	iter := c.db.Query(query, fmt.Sprintf("%d", clientID)).Iter()
	for iter.Scan(&pedidoIDStr, &dataPedido, &status, &valorTotal, &itensJSON) {
		pedidoID, err := strconv.ParseUint(pedidoIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter pedido_id: %v", err)
		}

		var itensMap []map[string]interface{}
		if err := json.Unmarshal([]byte(itensJSON), &itensMap); err != nil {
			return nil, fmt.Errorf("erro ao decodificar itens do pedido: %v", err)
		}

		var itens []model.OrderItem
		for _, itemMap := range itensMap {
			produtoIDStr := itemMap["produto_id"].(string)
			produtoID, err := strconv.ParseUint(produtoIDStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("erro ao converter produto_id: %v", err)
			}

			quantidade := int(itemMap["quantidade"].(float64))
			nomeProduto := itemMap["nome_produto"].(string)
			precoUnitario := itemMap["preco_unitario"].(float64)

			item := model.OrderItem{
				OrderID:   uint(pedidoID),
				ProductID: uint(produtoID),
				Quantity:  quantidade,
				Product: model.Product{
					ID:    uint(produtoID),
					Name:  nomeProduto,
					Price: precoUnitario,
				},
			}
			itens = append(itens, item)
		}

		order := model.Order{
			ID:         uint(pedidoID),
			ClientID:   clientID,
			OrderDate:  dataPedido,
			Status:     status,
			TotalValue: valorTotal,
			Itens:      itens,
		}
		orders = append(orders, order)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (c *CassandraRepository) Get5MostSoldProducts() ([]model.Product, error) {
	query := `SELECT id, nome, categoria, preco, estoque FROM produtos_por_vendas WHERE partition_key = 'all' LIMIT 5`
	var products []model.Product
	iter := c.db.Query(query).Iter()

	var (
		idStr     string
		nome      string
		categoria string
		preco     float64
		estoque   int
	)

	for iter.Scan(&idStr, &nome, &categoria, &preco, &estoque) {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter id do produto: %v", err)
		}

		product := model.Product{
			ID:       uint(id),
			Name:     nome,
			Category: categoria,
			Price:    preco,
			Stock:    estoque,
		}
		products = append(products, product)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return products, nil
}

func (c *CassandraRepository) GetLastMonthPixPayments() ([]model.Payment, error) {
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()
	mesAno := startDate.Format("2006-01")

	query := `
		SELECT id_pagamento, id_pedido, data_pagamento, valor_total
		FROM pagamentos_por_tipo_e_mes
		WHERE tipo = 'pix' AND mes_ano = ? AND data_pagamento >= ? AND data_pagamento <= ?
	`

	var payments []model.Payment
	iter := c.db.Query(query, mesAno, startDate, endDate).Iter()

	var (
		idPagamentoStr string
		idPedidoStr    string
		dataPagamento  time.Time
		valorTotal     float64
	)

	for iter.Scan(&idPagamentoStr, &idPedidoStr, &dataPagamento, &valorTotal) {
		idPagamento, err := strconv.ParseUint(idPagamentoStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter id_pagamento: %v", err)
		}

		idPedido, err := strconv.ParseUint(idPedidoStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter id_pedido: %v", err)
		}

		payment := model.Payment{
			ID:          uint(idPagamento),
			OrderID:     uint(idPedido),
			Type:        "pix",
			Status:      "aprovado",
			PaymentDate: dataPagamento,
		}
		payments = append(payments, payment)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (c *CassandraRepository) GetClientTotalSpentByPeriod(clientID uint, startDate time.Time, endDate time.Time) (float64, error) {
	mesAno := startDate.Format("2006-01")
	query := `
		SELECT SUM(valor_total)
		FROM pagamentos_por_tipo_e_mes
		WHERE id_cliente = ? AND mes_ano = ? AND data_pagamento >= ? AND data_pagamento <= ?
		ALLOW FILTERING
	`

	var total float64
	err := c.db.Query(query, fmt.Sprintf("%d", clientID), mesAno, startDate, endDate).Scan(&total)
	if err != nil && err == gocql.ErrNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return total, nil
}

func NewCassandraRepository() *CassandraRepository {
	config := config.LoadCassandraConfig()

	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Keyspace = config.Keyspace
	cluster.Consistency = config.Consistency
	cluster.ConnectTimeout = config.Timeout
	cluster.NumConns = config.PoolSize
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: config.Retries}

	session, err := cluster.CreateSession()
	if err != nil {
		panic(fmt.Sprintf("Erro ao conectar com o Cassandra: %v", err))
	}

	return &CassandraRepository{db: session}
}
