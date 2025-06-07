# TechMarket Showcase - Análise Comparativa de Bancos de Dados

Este projeto apresenta uma análise comparativa de desempenho entre PostgreSQL, MongoDB e Cassandra em um cenário de e-commerce.

## 📋 Índice

- [Instalação](#instalação)
- [Como Executar](#como-executar)
- [Modelagem e Decisões de Design](#modelagem-e-decisões-de-design)
- [Consultas Implementadas](#consultas-implementadas)
- [Análise de Performance](#análise-de-performance)

## 🚀 Instalação

### Pré-requisitos

- Docker e Docker Compose
- Node.js (v14+)
- Python (v3.8+)

### Configuração dos Bancos de Dados

```bash
# Clone o repositório
git clone https://github.com/seu-usuario/techmarket_showcase.git
cd techmarket_showcase

# Inicie os containers dos bancos de dados
docker-compose up -d
```

## 💻 Como Executar

```bash
# Instale as dependências
pip install -r requirements.txt

# Execute os benchmarks
python run_benchmarks.py
```

## 🎯 Modelagem e Decisões de Design

### PostgreSQL (Relacional)

```sql
CREATE TABLE clientes (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    nome VARCHAR(255),
    cpf VARCHAR(11)
);

CREATE TABLE produtos (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255),
    categoria VARCHAR(100),
    preco DECIMAL(10,2)
);

CREATE TABLE pedidos (
    id SERIAL PRIMARY KEY,
    cliente_id INTEGER REFERENCES clientes(id),
    data_pedido TIMESTAMP,
    status VARCHAR(50)
);

CREATE TABLE pagamentos (
    id SERIAL PRIMARY KEY,
    pedido_id INTEGER REFERENCES pedidos(id),
    valor DECIMAL(10,2),
    metodo VARCHAR(50),
    data_pagamento TIMESTAMP
);
```

### MongoDB (Documento)

```javascript
// Clientes
{
  _id: ObjectId,
  email: String,
  nome: String,
  cpf: String
}

// Produtos
{
  _id: ObjectId,
  nome: String,
  categoria: String,
  preco: Number
}

// Pedidos
{
  _id: ObjectId,
  cliente_id: ObjectId,
  data_pedido: Date,
  status: String,
  produtos: [
    {
      produto_id: ObjectId,
      quantidade: Number,
      preco_unitario: Number
    }
  ]
}

// Pagamentos
{
  _id: ObjectId,
  pedido_id: ObjectId,
  valor: Number,
  metodo: String,
  data_pagamento: Date
}
```

### Cassandra (Colunar)

```sql
CREATE TABLE clientes (
    id uuid PRIMARY KEY,
    email text,
    nome text,
    cpf text
);

CREATE TABLE produtos (
    id uuid,
    categoria text,
    nome text,
    preco decimal,
    PRIMARY KEY (categoria, id)
);

CREATE TABLE pedidos_por_cliente (
    cliente_id uuid,
    pedido_id uuid,
    data_pedido timestamp,
    status text,
    PRIMARY KEY (cliente_id, pedido_id)
);

CREATE TABLE pagamentos (
    id uuid,
    pedido_id uuid,
    valor decimal,
    metodo text,
    data_pagamento timestamp,
    PRIMARY KEY (id, data_pagamento)
);
```

## 📊 Consultas Implementadas

### 1. Busca de Cliente por Email

```sql
-- PostgreSQL
SELECT * FROM clientes WHERE email = ?;

-- MongoDB
db.clientes.find({ email: "?" });

-- Cassandra
SELECT * FROM clientes WHERE email = ?;
```

### 2. Produtos por Categoria

```sql
-- PostgreSQL
SELECT * FROM produtos WHERE categoria = ?;

-- MongoDB
db.produtos.find({ categoria: "?" });

-- Cassandra
SELECT * FROM produtos WHERE categoria = ?;
```

### 3. Produtos Entregues por Cliente

```sql
-- PostgreSQL
SELECT c.nome, p.*
FROM clientes c
JOIN pedidos p ON c.id = p.cliente_id
WHERE p.status = 'ENTREGUE';

-- MongoDB
db.pedidos.aggregate([
  { $match: { status: "ENTREGUE" } },
  { $lookup: { from: "clientes", localField: "cliente_id", foreignField: "_id", as: "cliente" } }
]);

-- Cassandra
SELECT * FROM pedidos_por_cliente WHERE status = 'ENTREGUE';
```

## 📈 Análise de Performance

### Operações de INSERT (média de registros/segundo)

| Banco de Dados | Cliente | Produto | Pedido | Pagamento |
| -------------- | ------- | ------- | ------ | --------- |
| PostgreSQL     | 48,000  | 78,000  | 35,000 | 36,000    |
| MongoDB        | 150,000 | 129,000 | 2,700  | 147,000   |
| Cassandra      | 18,000  | 24,000  | 9,000  | 31,000    |

### Operações de QUERY (média de registros/segundo)

| Banco de Dados | Cliente por Email | Produto por Categoria | Produtos Entregues | Top 5 Produtos |
| -------------- | ----------------- | --------------------- | ------------------ | -------------- |
| PostgreSQL     | 6,900,000         | 4,300,000             | 18,000,000         | 600,000        |
| MongoDB        | 2,000,000         | 1,800,000             | 30,000,000         | 550,000        |
| Cassandra      | 3,500,000         | 900,000               | 2,800,000          | 800,000        |

### Análise Crítica

1. **Operações de INSERT**:

   - MongoDB se destaca em inserções simples (Cliente, Produto, Pagamento)
   - PostgreSQL mantém consistência em todas as operações
   - Cassandra apresenta performance inferior, mas consistente
   - MongoDB tem queda significativa em inserções de Pedidos devido à complexidade do documento

2. **Operações de QUERY**:

   - PostgreSQL lidera em consultas simples por índice (email, categoria)
   - MongoDB se destaca em consultas que envolvem relacionamentos (produtos entregues)
   - Cassandra mantém performance estável, mas inferior em consultas complexas

3. **Recomendações**:

   - Para operações CRUD simples: MongoDB
   - Para consultas complexas com joins: PostgreSQL
   - Para escritas intensivas com schema flexível: MongoDB
   - Para alta disponibilidade e escalabilidade: Cassandra

4. **Considerações de Uso**:
   - PostgreSQL: Melhor para garantir ACID e consultas complexas
   - MongoDB: Ideal para dados não estruturados e alta velocidade de escrita
   - Cassandra: Recomendado para cenários de alta disponibilidade e distribuição global
