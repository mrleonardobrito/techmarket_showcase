package main

import (
	"log"
	"techmarket_showcase/benchmark"
	"techmarket_showcase/config"
	"techmarket_showcase/repo"
	"techmarket_showcase/repo/seed"
	"time"
)

const (
	CLIENT_INSERT_SIZE  = 20000
	PRODUCT_INSERT_SIZE = 5000
	ORDER_INSERT_SIZE   = 10000
	PAYMENT_INSERT_SIZE = 10000
)

func main() {
	config.LoadDotEnv()

	benchLogger, err := benchmark.NewBenchmarkLogger("benchmark_results.log")
	if err != nil {
		log.Fatalf("Erro ao criar benchmark logger: %v", err)
	}
	defer benchLogger.Close()

	mongoRepo := repo.NewMongoDBRepository()
	postgresRepo := repo.NewPostgresRepository()
	cassandraRepo := repo.NewCassandraRepository()

	clients := seed.GenerateClients(CLIENT_INSERT_SIZE)
	products := seed.GenerateProducts(PRODUCT_INSERT_SIZE)
	orders := seed.GenerateOrders(ORDER_INSERT_SIZE, CLIENT_INSERT_SIZE, PRODUCT_INSERT_SIZE)
	payments := seed.GeneratePayments(ORDER_INSERT_SIZE, PAYMENT_INSERT_SIZE)

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Insert, "Cliente", CLIENT_INSERT_SIZE, func() error {
		return postgresRepo.BatchCreateClient(clients)
	})

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Insert, "Produto", PRODUCT_INSERT_SIZE, func() error {
		return postgresRepo.BatchCreateProduct(products)
	})

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Insert, "Pedido", ORDER_INSERT_SIZE, func() error {
		return postgresRepo.BatchCreateOrder(orders)
	})

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Insert, "Pagamento", PAYMENT_INSERT_SIZE, func() error {
		return postgresRepo.BatchCreatePayment(payments)
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Insert, "Cliente", CLIENT_INSERT_SIZE, func() error {
		return mongoRepo.BatchCreateClient(clients)
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Insert, "Produto", PRODUCT_INSERT_SIZE, func() error {
		return mongoRepo.BatchCreateProduct(products)
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Insert, "Pedido", ORDER_INSERT_SIZE, func() error {
		return mongoRepo.BatchCreateOrder(orders)
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Insert, "Pagamento", PAYMENT_INSERT_SIZE, func() error {
		return mongoRepo.BatchCreatePayment(payments)
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Insert, "Cliente", CLIENT_INSERT_SIZE, func() error {
		return cassandraRepo.BatchCreateClient(clients)
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Insert, "Produto", PRODUCT_INSERT_SIZE, func() error {
		return cassandraRepo.BatchCreateProduct(products)
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Insert, "Pedido", ORDER_INSERT_SIZE, func() error {
		return cassandraRepo.BatchCreateOrder(orders)
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Insert, "Pagamento", PAYMENT_INSERT_SIZE, func() error {
		return cassandraRepo.BatchCreatePayment(payments)
	})

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Query, "Cliente por email", CLIENT_INSERT_SIZE, func() error {
		_, err := postgresRepo.GetClientByEmail("teste@teste.com")
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Query, "Cliente por email", CLIENT_INSERT_SIZE, func() error {
		_, err := mongoRepo.GetClientByEmail("teste@teste.com")
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Query, "Cliente por email", CLIENT_INSERT_SIZE, func() error {
		_, err := cassandraRepo.GetClientByEmail("teste@teste.com")
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Query, "Produto por categoria", PRODUCT_INSERT_SIZE, func() error {
		_, err := postgresRepo.GetProductByCategory("teste")
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Query, "Produto por categoria", PRODUCT_INSERT_SIZE, func() error {
		_, err := mongoRepo.GetProductByCategory("teste")
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Query, "Produto por categoria", PRODUCT_INSERT_SIZE, func() error {
		_, err := cassandraRepo.GetProductByCategory("teste")
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Query, "Produtos entregues por cliente", CLIENT_INSERT_SIZE, func() error {
		_, err := postgresRepo.GetDeliveredOrdersByClient(1)
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Query, "Produtos entregues por cliente", CLIENT_INSERT_SIZE, func() error {
		_, err := mongoRepo.GetDeliveredOrdersByClient(1)
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Query, "Produtos entregues por cliente", CLIENT_INSERT_SIZE, func() error {
		_, err := cassandraRepo.GetDeliveredOrdersByClient(1)
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Query, "5 produtos mais vendidos", PRODUCT_INSERT_SIZE, func() error {
		_, err := postgresRepo.Get5MostSoldProducts()
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Query, "5 produtos mais vendidos", PRODUCT_INSERT_SIZE, func() error {
		_, err := mongoRepo.Get5MostSoldProducts()
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Query, "5 produtos mais vendidos", PRODUCT_INSERT_SIZE, func() error {
		_, err := cassandraRepo.Get5MostSoldProducts()
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Query, "Pagamentos pix do último mês", PAYMENT_INSERT_SIZE, func() error {
		_, err := postgresRepo.GetLastMonthPixPayments()
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Query, "Pagamentos pix do último mês", PAYMENT_INSERT_SIZE, func() error {
		_, err := mongoRepo.GetLastMonthPixPayments()
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Query, "Pagamentos pix do último mês", PAYMENT_INSERT_SIZE, func() error {
		_, err := cassandraRepo.GetLastMonthPixPayments()
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Postgres, benchmark.Query, "Total gasto por cliente no último mês", CLIENT_INSERT_SIZE, func() error {
		_, err := postgresRepo.GetClientTotalSpentByPeriod(1, time.Now().AddDate(0, -1, 0), time.Now())
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.MongoDB, benchmark.Query, "Total gasto por cliente no último mês", CLIENT_INSERT_SIZE, func() error {
		_, err := mongoRepo.GetClientTotalSpentByPeriod(1, time.Now().AddDate(0, -1, 0), time.Now())
		if err != nil {
			return err
		}
		return nil
	})

	benchLogger.MeasureOperation(benchmark.Cassandra, benchmark.Query, "Total gasto por cliente no último mês", CLIENT_INSERT_SIZE, func() error {
		_, err := cassandraRepo.GetClientTotalSpentByPeriod(1, time.Now().AddDate(0, -1, 0), time.Now())
		if err != nil {
			return err
		}
		return nil
	})

}
