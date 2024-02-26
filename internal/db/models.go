package db

import (
	"sync"
	"time"
)

type Client struct {
	sync.Mutex

	ID      int `db:"id"`
	Balance int `db:"saldo"`
	Limit   int `db:"limite"`
}

type Transaction struct {
	ID          int       `db:"id"`
	Value       int       `db:"valor"`
	Type        string    `db:"tipo"`
	Description string    `db:"descricao"`
	CreatedAt   time.Time `db:"realizada_em"`
	ClientID    int       `db:"cliente_id"`
}

type AccountSummary struct {
	Balance       int       `json:"total"`
	Limit         int       `json:"limite"`
	StatementDate time.Time `json:"data_extrato"`
}

type TransactionSummary struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

type Extract struct {
	AccountSummary  AccountSummary       `json:"saldo"`
	LastTransations []TransactionSummary `json:"ultimas_transacoes"`
}
