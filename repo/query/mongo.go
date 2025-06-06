package repository

import (
	"context"
	"errors"
	"techmarket/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoRepository(client *mongo.Client, dbName string) QueryRepository {
	return &mongoRepository{
		client: client,
		db:     client.Database(dbName),
	}
}

func (r *mongoRepository) GetClientWithLastOrdersByEmail(ctx context.Context, email string) (model.Client, error) {
	var client model.Client
	collection := r.db.Collection("clientes")
	filter := bson.M{"email": email}

	err := collection.FindOne(ctx, filter).Decode(&client)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Client{}, errors.New("cliente não encontrado")
		}
		return model.Client{}, err
	}
	return client, nil
}

func (r *mongoRepository) GetProductByCategory(ctx context.Context, category string) ([]model.Product, error) {
	var products []model.Product
	collection := r.db.Collection("produtos")
	filter := bson.M{"categoria": category}
	opts := options.Find().SetSort(bson.D{{Key: "preco", Value: 1}}) 

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *mongoRepository) GetDeliveredProductsByClient(ctx context.Context, clientID uint) ([]model.Product, error) {
	var products []model.Product
	pedidosCollection := r.db.Collection("pedidos")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"id_cliente": clientID, "status": "entregue"}}},
		bson.D{{Key: "$unwind", Value: "$itens"}},
		bson.D{{Key: "$group", Value: bson.M{"_id": "$itens.produto_id"}}},
		bson.D{{Key: "$project", Value: bson.M{"produto_id": "$_id", "_id": 0}}},
	}

	cursor, err := pedidosCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []struct {
		ProdutoID uint `bson:"produto_id"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []model.Product{}, nil
	}

	var productIDs []uint
	for _, res := range results {
		productIDs = append(productIDs, res.ProdutoID)
	}

	produtosCollection := r.db.Collection("produtos")
	filter := bson.M{"_id": bson.M{"$in": productIDs}}
	prodCursor, err := produtosCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = prodCursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *mongoRepository) Get5MostSoldProducts(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	pedidosCollection := r.db.Collection("pedidos")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$unwind", Value: "$itens"}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":           "$itens.produto_id",
			"total_vendido": bson.M{"$sum": "$itens.quantidade"},
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"total_vendido": -1}}},
		bson.D{{Key: "$limit", Value: 5}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "produtos",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "detalhes_produto",
		}}},
		bson.D{{Key: "$replaceWith", Value: bson.M{"$arrayElemAt": []interface{}{"$detalhes_produto", 0}}}},
	}

	cursor, err := pedidosCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *mongoRepository) GetLastMonthPixPayments(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	collection := r.db.Collection("pagamentos")

	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	filter := bson.M{
		"tipo":           "pix",
		"data_pagamento": bson.M{"$gte": oneMonthAgo},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *mongoRepository) GetClientTotalSpentByPeriod(ctx context.Context, clientID uint, startDate time.Time, endDate time.Time) (float64, error) {
	pedidosCollection := r.db.Collection("pedidos")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"id_cliente": clientID,
			"data_pedido": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":        nil, 
			"totalGasto": bson.M{"$sum": "$valor_total"},
		}}},
	}

	cursor, err := pedidosCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	var results []struct {
		TotalGasto float64 `bson:"totalGasto"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, nil
	}

	return results[0].TotalGasto, nil
}