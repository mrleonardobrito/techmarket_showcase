CREATE TABLE cliente (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    telefone VARCHAR(20),
    data_cadastro TIMESTAMPTZ DEFAULT NOW(),
    cpf VARCHAR(14) UNIQUE NOT NULL
);

CREATE INDEX idx_cliente_email ON cliente (email);

CREATE TABLE produto (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    categoria VARCHAR(100) NOT NULL,
    preco DECIMAL(10, 2) NOT NULL,
    estoque INT NOT NULL
);

CREATE INDEX idx_produto_categoria ON produto (categoria);

CREATE TABLE pedido (
    id SERIAL PRIMARY KEY,
    id_cliente INT NOT NULL REFERENCES Cliente (id),
    data_pedido TIMESTAMPTZ DEFAULT NOW(),
    status VARCHAR(50) NOT NULL,
    valor_total DECIMAL(10, 2) NOT NULL
);

CREATE INDEX idx_pedido_id_cliente ON pedido (id_cliente);

CREATE INDEX idx_pedido_status ON pedido (status);

CREATE TABLE item_pedido (
    id_pedido INT NOT NULL REFERENCES pedido (id),
    id_produto INT NOT NULL REFERENCES produto (id),
    quantidade INT NOT NULL,
    preco_unitario DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (id_pedido, id_produto)
);

CREATE TABLE pagamento (
    id SERIAL PRIMARY KEY,
    id_pedido INT NOT NULL REFERENCES pedido (id),
    tipo VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    data_pagamento TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pagamento_id_pedido ON pagamento (id_pedido);

CREATE INDEX idx_pagamento_tipo ON pagamento (tipo);

CREATE INDEX idx_pagamento_data ON Pagamento (data_pagamento);