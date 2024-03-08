package db

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Client) ProcessTransaction(ctx context.Context, tx *pgx.Tx, value int, transaction string) (status int) {

	(*tx).QueryRow(ctx, transaction, value, c.ID).Scan(&c.Balance, &c.Limit)

	(*tx).Commit(ctx)
	return http.StatusOK

}

func (c *Client) ProcessCreditTransaction(ctx context.Context, pool *pgxpool.Pool, value int) (status int) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return http.StatusBadRequest
	}
	defer tx.Rollback(ctx)

	status = c.ProcessTransaction(ctx, &tx, value, credit_transaction)
	if status != http.StatusOK {
		return status
	}
	tx.Commit(ctx)
	return http.StatusOK
}

func (c *Client) ProcessDebitTransaction(ctx context.Context, pool *pgxpool.Pool, value int) (status int) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return http.StatusBadRequest
	}
	defer tx.Rollback(ctx)

	tx.QueryRow(ctx, "SELECT saldo, limite FROM clientes WHERE id = $1", c.ID).Scan(&c.Balance, &c.Limit)

	if c.Balance-value < -c.Limit {
		tx.Commit(ctx)
		return http.StatusUnprocessableEntity
	}

	status = c.ProcessTransaction(ctx, &tx, value, debit_transaction)
	if status != http.StatusOK {
		return status
	}
	tx.Commit(ctx)
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

func (t *Transaction) SaveToDB(ctx context.Context, pool *pgxpool.Pool) bool {

	_, err := pool.Exec(ctx, add_transaction, t.Value, t.Type, t.Description, t.ClientID)
	return err == nil
}

func GetExtrato(ctx context.Context, pool *pgxpool.Pool, c *Client) (*Extract, int) {

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, http.StatusBadRequest
	}
	defer tx.Rollback(ctx)

	extract := &Extract{
		AccountSummary:  AccountSummary{StatementDate: time.Now()},
		LastTransations: make([]TransactionSummary, 0, 10)}

	var Balance, Limit int

	row := tx.QueryRow(ctx, get_client_info, c.ID).Scan(&Balance, &Limit)
	if row != nil {
		return nil, http.StatusNotFound
	}

	extract.AccountSummary.Balance = Balance
	extract.AccountSummary.Limit = Limit

	rows, err := tx.Query(ctx, get_extract, c.ID)
	if err != nil {
		return nil, http.StatusBadRequest
	}
	defer rows.Close()

	for rows.Next() {
		var e TransactionSummary
		rows.Scan(&e.Value, &e.Type, &e.Description, &e.CreatedAt)
		extract.LastTransations = append(extract.LastTransations, e)
	}

	extract.AccountSummary.Balance = Balance
	extract.AccountSummary.Limit = Limit

	tx.Commit(ctx)

	return extract, http.StatusOK
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
