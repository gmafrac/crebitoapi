package routes

import (
	"context"

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
