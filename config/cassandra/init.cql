CREATE KEYSPACE IF NOT EXISTS techmarket
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

USE techmarket;

CREATE TABLE IF NOT EXISTS clientes_por_email (
    email text PRIMARY KEY,
    id uuid,
    nome text,
    telefone text,
    data_cadastro timestamp,
    cpf text
);

CREATE TABLE IF NOT EXISTS pedidos_por_cliente (
    id_cliente uuid,
    data_pedido timestamp,
    id_pedido uuid,
    status text,
    valor_total decimal,
    itens list<frozen<map<text, text>>>,
    PRIMARY KEY (id_cliente, data_pedido)
) WITH CLUSTERING ORDER BY (data_pedido DESC);

CREATE TABLE IF NOT EXISTS produtos_por_categoria (
    categoria text,
    preco decimal,
    id_produto uuid,
    nome text,
    estoque int,
    PRIMARY KEY (categoria, preco)
) WITH CLUSTERING ORDER BY (preco ASC);

CREATE TABLE IF NOT EXISTS vendas_de_produtos (
    id_produto uuid PRIMARY KEY,
    nome_produto text,
    total_vendido counter
);

CREATE TABLE IF NOT EXISTS pagamentos_por_tipo_e_mes (
    tipo text,
    mes_ano text,
    data_pagamento timestamp,
    id_pagamento uuid,
    id_pedido uuid,
    id_cliente uuid,
    PRIMARY KEY ((tipo, mes_ano), data_pagamento)
) WITH CLUSTERING ORDER BY (data_pagamento DESC);
