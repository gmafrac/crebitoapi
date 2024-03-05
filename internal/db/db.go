package db

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

func GetClient(pool *pgxpool.Pool, id int) (*Client, bool) {
	c := &Client{}

	err := pool.QueryRow(context.Background(), select_client, id).Scan(&c.ID, &c.Limit, &c.Balance)
	if err != nil {
		log.Print(err)
		return nil, false
	}
	return c, true
}

func (c *Client) ProcessCreditTransaction(pool *pgxpool.Pool, value int) {
	c.BalanceUpdate(pool, value)
}

func (c *Client) ProcessDebitTransaction(pool *pgxpool.Pool, value int) int {
	balance := c.Balance - value
	if balance < -c.Limit {
		return http.StatusUnprocessableEntity
	}
	return c.BalanceUpdate(pool, -value)

}

func (c *Client) BalanceUpdate(pool *pgxpool.Pool, value int) int {
	new_balance := c.Balance + value
	pool.Exec(
		context.Background(),
		update,
		new_balance, c.ID)
	c.Balance = new_balance

	return http.StatusOK

}

func (c *Client) SendResponse(w http.ResponseWriter) {

	json, err := json.Marshal(c)

	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func GetTransaction(r *http.Request, id int) (*Transaction, int) {
	transaction := &Transaction{}
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		return nil, http.StatusUnprocessableEntity
	}

	if len(transaction.Description) > 10 {
		return nil, http.StatusUnprocessableEntity
	}

	if transaction.Value == 0 || transaction.Type == "" || transaction.Description == "" {
		return nil, http.StatusUnprocessableEntity
	}

	defer r.Body.Close()
	transaction.ClientID = id

	return transaction, http.StatusOK
}

func (t *Transaction) SaveToDB(pool *pgxpool.Pool) bool {

	_, err := pool.Exec(
		context.Background(),
		add_transaction,
		t.Value, t.Type, t.Description, t.ClientID)
	return err == nil
}

func GetExtrato(pool *pgxpool.Pool, c *Client) (*Extract, bool) {

	extract := &Extract{
		AccountSummary: AccountSummary{
			Balance:       c.Balance,
			Limit:         c.Limit,
			StatementDate: time.Now()},
		LastTransations: []TransactionSummary{}}

	rows, err := pool.Query(
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
