package routes

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	// addr string
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewServer(pool *pgxpool.Pool, ctx context.Context) *Server {

	// addr := fmt.Sprintf("http://api0%s:%s", apiID, apiPort)
	// conn, err := pool.Acquire(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Release()

	// log.Printf("Database has been connected: " + utils.GetDBUrl())

	return &Server{
		// addr: addr,
		pool: pool,
		// conn: conn,
		ctx: ctx,
	}

}

// func (s *Server) GetAddress() string {
// 	return s.addr
// }

func (s *Server) IsAlive() bool {
	return true
}
