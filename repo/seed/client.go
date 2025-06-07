package seed

import (
	"fmt"
	"math/rand"
	"strings"
	"techmarket_showcase/model"
	"time"

	"github.com/go-faker/faker/v4"
)

func generateCPF() string {
	cpf := ""
	for range 9 {
		cpf += fmt.Sprintf("%d", rand.Intn(10))
	}

	sum := 0
	for i := range 9 {
		digit := int(cpf[i] - '0')
		sum += digit * (10 - i)
	}
	d1 := 11 - (sum % 11)
	if d1 >= 10 {
		d1 = 0
	}
	cpf += fmt.Sprintf("%d", d1)

	sum = 0
	for i := range 10 {
		digit := int(cpf[i] - '0')
		sum += digit * (11 - i)
	}
	d2 := 11 - (sum % 11)
	if d2 >= 10 {
		d2 = 0
	}
	cpf += fmt.Sprintf("%d", d2)

	return fmt.Sprintf("%s.%s.%s-%s",
		cpf[0:3],
		cpf[3:6],
		cpf[6:9],
		cpf[9:11])
}

func generatePhone() string {
	ddd := fmt.Sprintf("%02d", rand.Intn(89)+11)
	number := fmt.Sprintf("%09d", rand.Intn(100000000)+900000000)
	return fmt.Sprintf("%s%s", ddd, number)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateClients(count int) []model.Client {
	clients := make([]model.Client, count)

	for i := 0; i < count; i++ {
		var person struct {
			FirstName string `faker:"first_name"`
			LastName  string `faker:"last_name"`
			Domain    string `faker:"domain_name"`
		}
		faker.FakeData(&person)

		clients[i] = model.Client{
			Nome:      fmt.Sprintf("%s %s", person.FirstName, person.LastName),
			Email:     strings.ToLower(fmt.Sprintf("%s.%s@%s", person.FirstName, person.LastName, person.Domain)),
			Phone:     generatePhone(),
			CreatedAt: time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour), // Data aleatória no último ano
			CPF:       generateCPF(),
		}
	}
	return clients
}
