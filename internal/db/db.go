package db

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

const update = `
	UPDATE clientes
	SET saldo = $1
	WHERE id = $2
`
const add_transaction = `
	INSERT INTO transacoes (valor, tipo, descricao, cliente_id) 
	VALUES ($1, $2, $3, $4)
`
const select_client = `
	SELECT id, limite, saldo 
	FROM clientes
	WHERE clientes.id = $1
`
const get_extract = ` 
	SELECT valor, tipo, descricao, realizada_em
	FROM transacoes
	WHERE cliente_id = $1
	ORDER BY realizada_em DESC
	LIMIT 10	
`

func GetClient(conn *pgx.Conn, id int) (*Client, bool) {
	c := &Client{}
	c.Lock()
	defer c.Unlock()

	err := conn.QueryRow(context.Background(), select_client, id).Scan(&c.ID, &c.Limit, &c.Balance)
	if err != nil {
		return nil, false
	}
	return c, true
}

func (c *Client) ProcessCreditTransaction(conn *pgx.Conn, value int) {
	c.BalanceUpdate(conn, value)
}

func (c *Client) ProcessDebitTransaction(conn *pgx.Conn, value int) bool {
	balance := c.Balance - value
	if balance < c.Limit {
		return false
	}
	c.BalanceUpdate(conn, -value)
	return true
}

func (c *Client) BalanceUpdate(conn *pgx.Conn, value int) {

	conn.Exec(
		context.Background(),
		update,
		c.Balance+value, c.ID)
}

func GetTransaction(r *http.Request, id int) (*Transaction, bool) {
	transaction := &Transaction{}
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		return nil, false
	}
	defer r.Body.Close()
	transaction.ClientID = id
	return transaction, true
}

func (t *Transaction) SaveToDB(conn *pgx.Conn) bool {

	_, err := conn.Exec(
		context.Background(),
		add_transaction,
		t.Value, t.Type, t.Description, t.ClientID)
	return err == nil
}

func GetExtrato(conn *pgx.Conn, c *Client) (*Extract, bool) {

	extract := &Extract{
		AccountSummary: AccountSummary{
			Balance:       c.Balance,
			Limit:         c.Limit,
			StatementDate: time.Now()},
		LastTransations: []TransactionSummary{}}

	rows, err := conn.Query(
		context.Background(),
		get_extract, c.ID)

	if err != nil {
		return nil, false
	}

	for rows.Next() {
		var e TransactionSummary
		rows.Scan(&e.Value, &e.Type, &e.Description, &e.CreatedAt)
		extract.LastTransations = append(extract.LastTransations, e)
	}
	return extract, true
}

func (e *Extract) SendResponse(w http.ResponseWriter) {
	json, err := json.Marshal(e)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
