package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gmafrac/crebito_api/internal/db"
)

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	id, _ := strconv.Atoi(pathParts[1])

	cliente, ok := db.GetClient(s.pool, id)
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

		extract, ok := db.GetExtrato(s.pool, cliente)
		if !ok {
			http.Error(w, "Error", http.StatusBadRequest)
			return
		}

		extract.SendResponse(w)

		return

	case http.MethodPost:

		if pathParts[2] != "transacoes" {
			http.Error(w, "Error: Path not found", http.StatusUnprocessableEntity)
			return
		}

		transaction, status := db.GetTransaction(r, id)
		if status != http.StatusOK {
			http.Error(w, "Error decoding Json", status)
			return
		}

		transaction.SaveToDB(s.pool)

		switch transaction.Type {

		case "d":
			status := cliente.ProcessDebitTransaction(s.pool, transaction.Value)
			if status != http.StatusOK {
				http.Error(w, "", status)
			}
			cliente.SendResponse(w)
			return

		case "c":
			cliente.ProcessCreditTransaction(s.pool, transaction.Value)
			cliente.SendResponse(w)
			return

		default:
			http.Error(w, "Invalid transaction type", http.StatusUnprocessableEntity)
			return
		}

	default:
		http.Error(w, "Expect method GET or POST at /clientes/", http.StatusBadRequest)
	}

}
