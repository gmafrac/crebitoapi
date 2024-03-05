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

	pool, err := pgxpool.New(ctx, utils.GetDBUrl())
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	apiServer := routes.NewServer(pool, ctx)

	http.HandleFunc("/clientes/", apiServer.Handler)

	log.Printf("Server started")

	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Failed to start")
	}

}
