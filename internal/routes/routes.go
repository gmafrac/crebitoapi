package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gmafrac/crebito_api/internal/db"
	"github.com/jackc/pgx/v5"
)

type Server struct {
	conn *pgx.Conn
}

func NewServer(conn *pgx.Conn) *Server {
	return &Server{conn: conn}
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	id, _ := strconv.Atoi(pathParts[1])

	cliente, ok := db.GetClient(s.conn, id)
	if !ok {
		http.Error(w, "Error: Invalid client id", http.StatusNotFound)
		return
	}

	switch r.Method {

	case http.MethodGet:

		if pathParts[2] != "extrato" {
			http.Error(w, "Error: Path not found", http.StatusNotFound)
			return
		}

		extract, ok := db.GetExtrato(s.conn, cliente)
		if !ok {
			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		extract.SendResponse(w)

		return

	case http.MethodPost:

		if pathParts[2] != "transacoes" {
			http.Error(w, "Error: Path not found", http.StatusNotFound)
			return
		}

		transaction, ok := db.GetTransaction(r, id)
		if !ok {
			http.Error(w, "Error decoding Json", http.StatusBadRequest)
			return
		}

		transaction.SaveToDB(s.conn)

		switch transaction.Type {

		case "d":
			ok = cliente.ProcessDebitTransaction(s.conn, transaction.Value)
			if !ok {
				http.Error(w, "Insufficient Balance", http.StatusUnprocessableEntity)
			}
			return

		case "c":
			cliente.ProcessCreditTransaction(s.conn, transaction.Value)
			return

		default:
			http.Error(w, "Invalid transaction type", http.StatusUnprocessableEntity)
			return
		}

	default:
		http.Error(w, "Expect method GET or POST at /clientes/", http.StatusBadRequest)
	}

}
