db = db.getSiblingDB('techmarket_db');

db = db.getSiblingDB('admin');
db.auth('root', 'root_password');

db.createCollection('produtos');
db.produtos.createIndex({ categoria: 1, preco: 1 });

db.createCollection('pagamentos');
db.pagamentos.createIndex({ tipo: 1, data_pagamento: -1 });

db.createCollection('clientes');
db.clientes.createIndex({ email: 1 }, { unique: true });

db.clientes.insertOne({
    nome: "Maria Oliveira",
    email: "maria.o@example.com",
    telefone: "21988776655",
    data_cadastro: new ISODate(),
    cpf: "98765432109",
    pedidos: [
        {
            pedido_id: new ObjectId(),
            data_pedido: new ISODate("2024-10-25T14:30:00Z"),
            status: "entregue",
            valor_total: 4500.00,
            itens: [
                { produto_id: new ObjectId(), nome_produto: "Notebook Pro", quantidade: 1, preco_unitario: 4500.00 }
            ],
            pagamento: {
                pagamento_id: new ObjectId(),
                tipo: "cartao",
                status: "aprovado"
            }
        }
    ]
});

print("Banco de dados 'techmarket_db' e coleções criadas com sucesso!");
