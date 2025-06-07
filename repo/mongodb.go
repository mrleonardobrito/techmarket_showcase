package repo

import (
	"context"
	"fmt"
	"techmarket_showcase/config"
	"techmarket_showcase/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ TechMarketRepository = &MongoDBRepository{}

type MongoDBRepository struct {
	db *mongo.Client
}

func NewMongoDBRepository() *MongoDBRepository {
	config := config.LoadMongoDBConfig()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.URI))
	if err != nil {
		panic(fmt.Sprintf("Erro ao conectar com o MongoDB: %v", err))
	}

	return &MongoDBRepository{db: client}
}

func (m *MongoDBRepository) BatchCreateClient(clients []model.Client) error {
	collection := m.db.Database("techmarket_db").Collection("clientes")
	ctx := context.Background()

	var documents []any
	for _, client := range clients {
		doc := bson.M{
			"nome":          client.Nome,
			"email":         client.Email,
			"telefone":      client.Phone,
			"data_cadastro": client.CreatedAt,
			"cpf":           client.CPF,
			"pedidos":       []any{},
		}
		documents = append(documents, doc)
	}

	_, err := collection.InsertMany(ctx, documents)
	return err
}

func (m *MongoDBRepository) BatchCreateProduct(products []model.Product) error {
	collection := m.db.Database("techmarket_db").Collection("produtos")
	ctx := context.Background()

	var documents []any
	for _, product := range products {
		doc := bson.M{
			"nome":      product.Name,
			"categoria": product.Category,
			"preco":     product.Price,
			"estoque":   product.Stock,
		}
		documents = append(documents, doc)
	}

	_, err := collection.InsertMany(ctx, documents)
	return err
}

func (m *MongoDBRepository) BatchCreateOrder(orders []model.Order) error {
	clientsCollection := m.db.Database("techmarket_db").Collection("clientes")
	ctx := context.Background()

	for _, order := range orders {
		var itens []bson.M
		for _, item := range order.Itens {
			itemDoc := bson.M{
				"produto_id":     item.ProductID,
				"nome_produto":   item.Product.Name,
				"quantidade":     item.Quantity,
				"preco_unitario": item.Product.Price,
			}
			itens = append(itens, itemDoc)
		}

		pedidoDoc := bson.M{
			"pedido_id":   order.ID,
			"data_pedido": order.OrderDate,
			"status":      order.Status,
			"valor_total": order.TotalValue,
			"itens":       itens,
		}

		filter := bson.M{"_id": order.ClientID}
		update := bson.M{
			"$push": bson.M{
				"pedidos": pedidoDoc,
			},
		}

		_, err := clientsCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MongoDBRepository) BatchCreateOrderItem(orderItems []model.OrderItem) error {
	// No MongoDB, os itens do pedido são armazenados diretamente no documento do pedido
	// dentro da coleção de clientes, então esta função não precisa fazer nada
	return nil
}

func (m *MongoDBRepository) BatchCreatePayment(payments []model.Payment) error {
	collection := m.db.Database("techmarket_db").Collection("pagamentos")
	ctx := context.Background()

	var documents []any
	for _, payment := range payments {
		doc := bson.M{
			"tipo":           payment.Type,
			"status":         payment.Status,
			"data_pagamento": payment.PaymentDate,
			"pedido_id":      payment.OrderID,
		}
		documents = append(documents, doc)
	}

	_, err := collection.InsertMany(ctx, documents)
	return err
}

func (m *MongoDBRepository) GetClientByEmail(email string) (model.Client, error) {
	collection := m.db.Database("techmarket_db").Collection("clientes")
	ctx := context.Background()

	filter := bson.M{"email": email}
	var client model.Client
	err := collection.FindOne(ctx, filter).Decode(&client)
	if err != nil && err == mongo.ErrNoDocuments {
		return model.Client{}, nil
	}
	if err != nil {
		return model.Client{}, err
	}
	return client, nil
}

func (m *MongoDBRepository) GetDeliveredOrdersByClient(clientID uint) ([]model.Order, error) {
	collection := m.db.Database("techmarket_db").Collection("clientes")
	ctx := context.Background()

	filter := bson.M{"_id": clientID}
	var result struct {
		Pedidos []struct {
			PedidoID   uint      `bson:"pedido_id"`
			DataPedido time.Time `bson:"data_pedido"`
			Status     string    `bson:"status"`
			ValorTotal float64   `bson:"valor_total"`
			Itens      []struct {
				NomeProduto   string  `bson:"nome_produto"`
				ProdutoID     uint    `bson:"produto_id"`
				Quantidade    int     `bson:"quantidade"`
				PrecoUnitario float64 `bson:"preco_unitario"`
			} `bson:"itens"`
		} `bson:"pedidos"`
	}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var orders []model.Order
	for _, pedido := range result.Pedidos {
		if pedido.Status == "entregue" {
			var itens []model.OrderItem
			for _, item := range pedido.Itens {
				itens = append(itens, model.OrderItem{
					OrderID:   pedido.PedidoID,
					ProductID: item.ProdutoID,
					Quantity:  item.Quantidade,
					Product: model.Product{
						ID:    item.ProdutoID,
						Name:  item.NomeProduto,
						Price: item.PrecoUnitario,
					},
				})
			}

			order := model.Order{
				ID:         pedido.PedidoID,
				OrderDate:  pedido.DataPedido,
				Status:     pedido.Status,
				TotalValue: pedido.ValorTotal,
				ClientID:   clientID,
				Itens:      itens,
			}
			orders = append(orders, order)
		}
	}

	if len(orders) == 0 {
		return nil, nil
	}

	return orders, nil
}

func (m *MongoDBRepository) Get5MostSoldProducts() ([]model.Product, error) {
	collection := m.db.Database("techmarket_db").Collection("produtos")
	ctx := context.Background()

	filter := bson.M{}
	var products []model.Product
	cursor, err := collection.Find(ctx, filter)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (m *MongoDBRepository) GetProductByCategory(category string) ([]model.Product, error) {
	collection := m.db.Database("techmarket_db").Collection("produtos")
	ctx := context.Background()

	filter := bson.M{"categoria": category}
	var products []model.Product
	cursor, err := collection.Find(ctx, filter)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (m *MongoDBRepository) GetLastMonthPixPayments() ([]model.Payment, error) {
	collection := m.db.Database("techmarket_db").Collection("pagamentos")
	ctx := context.Background()

	filter := bson.M{"tipo": "pix", "data_pagamento": bson.M{"$gte": time.Now().AddDate(0, -1, 0), "$lte": time.Now()}}
	var payments []model.Payment
	cursor, err := collection.Find(ctx, filter)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &payments)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (m *MongoDBRepository) GetClientTotalSpentByPeriod(clientID uint, startDate time.Time, endDate time.Time) (float64, error) {
	collection := m.db.Database("techmarket_db").Collection("pagamentos")
	ctx := context.Background()

	filter := bson.M{"id_cliente": clientID, "data_pagamento": bson.M{"$gte": startDate, "$lte": endDate}}
	var payments []model.Payment
	cursor, err := collection.Find(ctx, filter)
	if err != nil && err == mongo.ErrNoDocuments {
		return 0, nil
	}

	err = cursor.All(ctx, &payments)
	if err != nil {
		return 0, err
	}

	var total float64
	collection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{"id_cliente": clientID, "data_pagamento": bson.M{"$gte": startDate, "$lte": endDate}},
		},
		{
			"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$valor_total"}},
		},
	})
	return total, nil
}
