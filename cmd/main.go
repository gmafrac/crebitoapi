package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gmafrac/crebito_api/internal/routes"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	mux := http.NewServeMux()
	server := routes.NewServer(conn)

	mux.HandleFunc("/clientes/", server.Handler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
