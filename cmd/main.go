package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gmafrac/crebito_api/internal/routes"
	"github.com/gmafrac/crebito_api/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	utils.LoadEnv()

	ctx := context.Background()

	log.Print("Starting server...", utils.GetDBUrl())
	pool, err := pgxpool.New(ctx, utils.GetDBUrl())
	if err != nil {
		log.Print("Failed to connect to database")
		panic(err)
	}
	defer pool.Close()

	apiServer := routes.NewServer(pool, ctx)

	http.HandleFunc("/clientes/", apiServer.Handler)

	log.Print("Server started at :8000")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatalf("Failed to start")
	}

}
