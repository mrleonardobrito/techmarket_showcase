package seed

import (
	"fmt"
	"math/rand"
	"techmarket_showcase/model"
)

var (
	categories = []string{
		"Smartphones",
		"Notebooks",
		"Tablets",
		"Smart TVs",
		"Fones de Ouvido",
		"Smartwatches",
		"Câmeras",
		"Acessórios",
		"Periféricos",
		"Componentes PC",
	}

	productPrefixes = []string{
		"Pro",
		"Ultra",
		"Max",
		"Plus",
		"Lite",
		"Premium",
		"Elite",
		"Smart",
		"Tech",
		"Advanced",
	}

	productBrands = []string{
		"TechPro",
		"SmartTech",
		"InnovatePro",
		"FutureTech",
		"NextGen",
		"EliteTech",
		"PrimeTech",
		"UltraTech",
		"MaxTech",
		"ProTech",
	}

	categoryPriceRanges = map[string]struct{ min, max float64 }{
		"Smartphones":     {800.00, 8000.00},
		"Notebooks":       {2000.00, 15000.00},
		"Tablets":         {500.00, 5000.00},
		"Smart TVs":       {1200.00, 12000.00},
		"Fones de Ouvido": {50.00, 2000.00},
		"Smartwatches":    {200.00, 3000.00},
		"Câmeras":         {500.00, 8000.00},
		"Acessórios":      {20.00, 500.00},
		"Periféricos":     {50.00, 1000.00},
		"Componentes PC":  {100.00, 5000.00},
	}
)

func generateProductName(category string) string {
	brand := productBrands[rand.Intn(len(productBrands))]
	prefix := productPrefixes[rand.Intn(len(productPrefixes))]
	model := fmt.Sprintf("%d", 1000+rand.Intn(9000))

	return fmt.Sprintf("%s %s %s %s", brand, category, prefix, model)
}

func generatePrice(category string) float64 {
	priceRange := categoryPriceRanges[category]
	price := priceRange.min + rand.Float64()*(priceRange.max-priceRange.min)
	return float64(int(price*100)) / 100
}

func GenerateProducts(count int) []model.Product {
	products := make([]model.Product, count)

	for i := range count {
		category := categories[rand.Intn(len(categories))]

		products[i] = model.Product{
			Name:     generateProductName(category),
			Category: category,
			Price:    generatePrice(category),
			Stock:    rand.Intn(901) + 100,
		}
	}

	return products
}
