package routes

import (
	"context"
	"net/http"

	"github.com/gmafrac/crebito_api/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewServer(pool *pgxpool.Pool, ctx context.Context) *Server {
	return &Server{
		pool: pool,
		ctx:  ctx,
	}
}

func (s *Server) Transaction(w http.ResponseWriter, r *http.Request, client *db.Client) {
	ctx := s.ctx

	transaction, status := db.GetTransaction(r, client.ID)
	if status != http.StatusOK {
		http.Error(w, "", status)
		return
	}

	switch transaction.Type {

	case "d":
		status := client.ProcessDebitTransaction(ctx, s.pool, transaction.Value)
		if status != http.StatusOK {
			http.Error(w, "", status)
			return
		}

	case "c":
		client.ProcessCreditTransaction(ctx, s.pool, transaction.Value)

	default:
		http.Error(w, "", http.StatusUnprocessableEntity)
	}

	transaction.SaveToDB(ctx, s.pool)
	client.SendResponse(w)

}

func (s *Server) Extract(w http.ResponseWriter, client *db.Client) {

	extract, status := db.GetExtrato(s.ctx, s.pool, client)
	if status != http.StatusOK {
		http.Error(w, "", status)
		return
	}

	extract.SendResponse(w)
}
