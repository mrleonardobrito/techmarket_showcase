CREATE KEYSPACE IF NOT EXISTS techmarket
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

USE techmarket;

DROP TABLE IF EXISTS clientes_por_email;
DROP TABLE IF EXISTS pedidos_por_cliente;
DROP TABLE IF EXISTS produtos_por_categoria;
DROP TABLE IF EXISTS pagamentos_por_tipo_e_mes;

CREATE TABLE IF NOT EXISTS clientes_por_email (
    email text PRIMARY KEY,
    id text,
    nome text,
    telefone text,
    data_cadastro timestamp,
    cpf text
);

CREATE TABLE IF NOT EXISTS pedidos_por_cliente (
    id_cliente text,
    pedido_id text,
    data_pedido timestamp,
    status text,
    valor_total decimal,
    itens text,
    PRIMARY KEY (id_cliente, pedido_id)
);

CREATE INDEX IF NOT EXISTS idx_pedidos_status ON pedidos_por_cliente (status);

CREATE TABLE IF NOT EXISTS produtos_por_categoria (
    categoria text,
    preco double,
    id_produto text,
    nome text,
    estoque int,
    PRIMARY KEY (categoria, preco)
) WITH CLUSTERING ORDER BY (preco ASC);

CREATE TABLE IF NOT EXISTS produtos_vendas_counter (
    id_produto uuid,
    nome_produto text,
    PRIMARY KEY (id_produto)
);

CREATE TABLE IF NOT EXISTS produtos_total_vendido (
    id_produto uuid PRIMARY KEY,
    total_vendido counter
);

CREATE TABLE IF NOT EXISTS pagamentos_por_tipo_e_mes (
    tipo text,
    mes_ano text,
    data_pagamento timestamp,
    id_pagamento text,
    id_pedido text,
    id_cliente text,
    valor_total decimal,
    PRIMARY KEY ((tipo, mes_ano), data_pagamento)
) WITH CLUSTERING ORDER BY (data_pagamento DESC);

CREATE TABLE IF NOT EXISTS produtos_por_vendas (
    partition_key text,
    id text,
    nome text,
    categoria text,
    preco decimal,
    estoque int,
    total_vendas counter,
    PRIMARY KEY (partition_key, total_vendas, id)
) WITH CLUSTERING ORDER BY (total_vendas DESC);

-- Inserir um registro inicial para a partition key 'all'
INSERT INTO produtos_por_vendas (partition_key, id, nome, categoria, preco, estoque) 
VALUES ('all', '1', 'Produto Inicial', 'Categoria', 0.0, 0);
