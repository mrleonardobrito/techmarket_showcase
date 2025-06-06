import (
	"context"
	"database/sql"
	"errors"
	"techmarket/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryRepository interface {
	GetClientByEmail(ctx context.Context, email string) (model.Client, error)
	GetProductByCategory(ctx context.Context, category string) ([]model.Product, error)
	GetDeliveredProductsByClient(ctx context.Context, clientID uint) ([]model.Product, error)
	Get5MostSoldProducts(ctx context.Context) ([]model.Product, error)
	GetLastMonthPixPayments(ctx context.Context) ([]model.Payment, error)
	GetClientTotalSpentByPeriod(ctx context.Context, clientID uint, startDate time.Time, endDate time.Time) (float64, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(dbpool *pgxpool.Pool) QueryRepository {
	return &postgresRepository{db: dbpool}
}

func (r *postgresRepository) GetClientWithLastOrdersByEmail(ctx context.Context, email string) (model.Client, []model.Pedido, error) {
	query := `
		SELECT 
			c.id, c.nome, c.email, c.telefone, c.data_cadastro, c.cpf,
			p.id, p.data_pedido, p.total
		FROM Cliente c
		LEFT JOIN Pedido p ON c.id = p.id_cliente
		WHERE c.email = $1
		ORDER BY p.data_pedido DESC
		LIMIT 3;
	`

	rows, err := r.db.Query(ctx, query, email)
	if err != nil {
		return model.Client{}, nil, err
	}
	defer rows.Close()

	var client model.Client
	var pedidos []model.Pedido
	clientScanned := false

	for rows.Next() {
		var (
			pedido model.Pedido
			pedidoID *int
			dataPedido *time.Time
			total *float64
		)

		if !clientScanned {
			err := rows.Scan(
				&client.ID, &client.Nome, &client.Email, &client.Telefone, &client.DataCadastro, &client.CPF,
				&pedidoID, &dataPedido, &total,
			)
			if err != nil {
				return model.Client{}, nil, err
			}
			clientScanned = true
		} else {
			err := rows.Scan(
				new(interface{}), new(interface{}), new(interface{}), new(interface{}), new(interface{}), new(interface{}),
				&pedidoID, &dataPedido, &total,
			)
			if err != nil {
				return model.Client{}, nil, err
			}
		}

		if pedidoID != nil {
			pedido.ID = *pedidoID
			pedido.ClienteID = client.ID
			pedido.DataPedido = *dataPedido
			pedido.Total = *total
			pedidos = append(pedidos, pedido)
		}
	}

	if !clientScanned {
		return model.Client{}, nil, errors.New("cliente não encontrado")
	}

	return client, pedidos, nil
}



func (r *postgresRepository) GetProductByCategory(ctx context.Context, category string) ([]model.Product, error) {
	var products []model.Product
	query := ` SELECT id, nome, categoria, preco, estoque FROM Produto WHERE categoria = $1 ORDER BY preco ASC `

	rows, err := r.db.Query(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Nome, &p.Categoria, &p.Preco, &p.Estoque); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *postgresRepository) GetDeliveredProductsByClient(ctx context.Context, clientID uint) ([]model.Product, error) {
	var products []model.Product
	query := `
		SELECT DISTINCT p.id, p.nome, p.categoria, p.preco, p.estoque
		FROM Produto p
		JOIN Item_Pedido ip ON p.id = ip.id_produto
		JOIN Pedido o ON ip.id_pedido = o.id
		WHERE o.id_cliente = $1 AND o.status = 'entregue' `

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Nome, &p.Categoria, &p.Preco, &p.Estoque); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *postgresRepository) Get5MostSoldProducts(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	query := `
		SELECT p.id, p.nome, p.categoria, p.preco, p.estoque
		FROM Produto p
		JOIN Item_Pedido ip ON p.id = ip.id_produto
		GROUP BY p.id
		ORDER BY SUM(ip.quantidade) DESC
		LIMIT 5 `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Nome, &p.Categoria, &p.Preco, &p.Estoque); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *postgresRepository) GetLastMonthPixPayments(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	query := `
		SELECT id, id_pedido, tipo, status, data_pagamento 
		FROM Pagamento 
		WHERE tipo = 'pix' AND data_pagamento >= NOW() - INTERVAL '1 month' `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(&p.ID, &p.IDPedido, &p.Tipo, &p.Status, &p.DataPagamento); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (r *postgresRepository) GetClientTotalSpentByPeriod(ctx context.Context, clientID uint, startDate time.Time, endDate time.Time) (float64, error) {
	var totalSpent float64
	query := `
		SELECT COALESCE(SUM(valor_total), 0) 
		FROM Pedido 
		WHERE id_cliente = $1 AND data_pedido BETWEEN $2 AND $3 `

	err := r.db.QueryRow(ctx, query, clientID, startDate, endDate).Scan(&totalSpent)
	if err != nil {
		return 0, err
	}

	return totalSpent, nil
}


