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

	lb := routes.NewLoadBalancer(
		utils.GetEnv("SERVER_ADDR"),
		[]*routes.Server{
			routes.NewServer(conn, "http://localhost:8081"),
			routes.NewServer(conn, "http://localhost:8082"),
		})

	http.HandleFunc("/clientes/", lb.ServeProxy)

	log.Printf("Server started at port " + lb.GetPort())
	if err := http.ListenAndServe(":"+lb.GetPort(), nil); err != nil {
		log.Fatalf("Failed to start")
	}
}
