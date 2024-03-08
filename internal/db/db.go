package db

const add_transaction = `
	INSERT INTO transacoes (valor, tipo, descricao, cliente_id) 
	VALUES ($1, $2, $3, $4)
`

const credit_transaction = `
	UPDATE clientes
	SET saldo = saldo + $1
	WHERE id = $2
	RETURNING saldo, limite
`

const debit_transaction = `
	UPDATE clientes
	SET saldo = saldo - $1
	WHERE id = $2
	RETURNING saldo, limite
`

const get_client_info = `
	SELECT SALDO, LIMITE 
	FROM CLIENTES 
	WHERE ID = $1;
`

const get_extract = ` 
	SELECT t.valor, t.tipo, t.descricao, t.realizada_em
	FROM transacoes t
	WHERE t.cliente_id = $1
	ORDER BY t.realizada_em DESC
	LIMIT 10;
`
