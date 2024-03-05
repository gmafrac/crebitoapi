package db

import (
	"time"
)

type Client struct {
	ID      int `db:"id" json:"-"`
	Balance int `db:"saldo" json:"saldo"`
	Limit   int `db:"limite" json:"limite"`
}

type Transaction struct {
	ID          int       `db:"id" json:"id"`
	Value       int       `db:"valor" json:"valor"`
	Type        string    `db:"tipo" json:"tipo"`
	Description string    `db:"descricao" json:"descricao"`
	CreatedAt   time.Time `db:"realizada_em" json:"realizada_em"`
	ClientID    int       `db:"cliente_id" json:"cliente_id"`
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
