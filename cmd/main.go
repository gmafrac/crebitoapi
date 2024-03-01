package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gmafrac/crebito_api/internal/routes"
	"github.com/gmafrac/crebito_api/internal/utils"
	"github.com/jackc/pgx/v5"
)

func main() {
	utils.LoadEnv()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, utils.GetDBUrl())
	if err != nil {
		log.Print("Error connecting to the database")
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	apiServer := routes.NewServer(conn, utils.GetEnv("API_ID"))

	http.HandleFunc("/clientes/", apiServer.Handler)

	log.Printf("Server started at port: " + apiServer.GetAddress())

	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatalf("Failed to start")
	}
}
