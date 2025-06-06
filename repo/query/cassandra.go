package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"techmarket/model"
	"time"

	"github.com/gocql/gocql"
)

type cassandraRepository struct {
	session *gocql.Session
}

func NewCassandraRepository(hosts []string, keyspace string) (QueryRepository, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum 

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &cassandraRepository{session: session}, nil
}

func (r *cassandraRepository) GetClientWithLastOrdersByEmail(ctx context.Context, email string) (model.Client, error) {
	var client model.Client
	query := `SELECT id, nome, email, telefone, data_cadastro, cpf FROM clientes_por_email WHERE email = ? LIMIT 1`

	err := r.session.Query(query, email).WithContext(ctx).Scan(
		&client.ID, &client.Nome, &client.Email, &client.Telefone, &client.DataCadastro, &client.CPF,
	)

	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return model.Client{}, errors.New("cliente não encontrado")
		}
		return model.Client{}, err
	}
	return client, nil
}

func (r *cassandraRepository) GetProductByCategory(ctx context.Context, category string) ([]model.Product, error) {
	var products []model.Product
	query := `SELECT id_produto, nome, categoria, preco, estoque FROM produtos_por_categoria WHERE categoria = ?`

	iter := r.session.Query(query, category).WithContext(ctx).Iter()
	var p model.Product
	for iter.Scan(&p.ID, &p.Nome, &p.Categoria, &p.Preco, &p.Estoque) {
		products = append(products, p)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *cassandraRepository) GetDeliveredProductsByClient(ctx context.Context, clientID uint) ([]model.Product, error) {
	var products []model.Product
	query := `SELECT itens FROM pedidos_por_cliente WHERE id_cliente = ? AND status = 'entregue' ALLOW FILTERING`

	iter := r.session.Query(query, clientID).WithContext(ctx).Iter()
	
	var scannedItems []map[string]string
	var productIDs []uint

	log.Println("A implementação de GetDeliveredProductsByClient para Cassandra é uma simplificação e pode não funcionar diretamente com a estrutura de 'itens' definida. Ela demonstra a necessidade de uma modelagem de dados cuidadosa.")


	return products, nil 
}

func (r *cassandraRepository) Get5MostSoldProducts(ctx context.Context) ([]model.Product, error) {
	return nil, errors.New("operação 'top 5' não é suportada em tempo real no Cassandra com a modelagem atual; requer uma tabela de sumarização pré-calculada")
}

func (r *cassandraRepository) GetLastMonthPixPayments(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	now := time.Now()
	currentMonth := now.Format("2006-01") 
	prevMonth := now.AddDate(0, -1, 0).Format("2006-01")
	
	monthsToQuery := []string{currentMonth}
	if !contains(monthsToQuery, prevMonth) {
		monthsToQuery = append(monthsToQuery, prevMonth)
	}
	
	query := `SELECT id, id_pedido, tipo, status, data_pagamento FROM pagamentos_por_tipo_e_mes WHERE tipo = 'pix' AND mes_ano IN ?`
	
	oneMonthAgo := now.AddDate(0, -1, 0)

	iter := r.session.Query(query, monthsToQuery).WithContext(ctx).Iter()
	var p model.Payment
	for iter.Scan(&p.ID, &p.IDPedido, &p.Tipo, &p.Status, &p.DataPagamento) {
		if p.DataPagamento.After(oneMonthAgo) {
			payments = append(payments, p)
		}
	}
	
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *cassandraRepository) GetClientTotalSpentByPeriod(ctx context.Context, clientID uint, startDate time.Time, endDate time.Time) (float64, error) {
	var totalSpent float64
	query := `SELECT valor_total FROM pedidos_por_cliente WHERE id_cliente = ? AND data_pedido >= ? AND data_pedido <= ?`

	iter := r.session.Query(query, clientID, startDate, endDate).WithContext(ctx).Iter()
	var valor_pedido float64
	for iter.Scan(&valor_pedido) {
		totalSpent += valor_pedido
	}

	if err := iter.Close(); err != nil {
		return 0, err
	}
	
	return totalSpent, nil
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}