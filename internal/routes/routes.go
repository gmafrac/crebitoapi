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

	cliente, status := db.GetClient(s.pool, id)
	if status != http.StatusOK {
		http.Error(w, "", status)
		return
	}

	switch r.Method {

	case http.MethodGet:

		if pathParts[2] != "extrato" {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		extract, status := db.GetExtrato(s.pool, cliente)
		if status != http.StatusOK {
			http.Error(w, "", status)
			return
		}

		extract.SendResponse(w)

	case http.MethodPost:

		if pathParts[2] != "transacoes" {
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}

		transaction, status := db.GetTransaction(r, id)
		if status != http.StatusOK {
			http.Error(w, "", status)
			return
		}

		transaction.SaveToDB(s.pool)

		switch transaction.Type {

		case "d":
			status := cliente.ProcessDebitTransaction(s.pool, transaction.Value)

			if status != http.StatusOK {
				http.Error(w, "", status)
				return
			}
			cliente.SendResponse(w)

		case "c":
			cliente.ProcessCreditTransaction(s.pool, transaction.Value)
			cliente.SendResponse(w)

		default:
			http.Error(w, "", http.StatusUnprocessableEntity)
		}

	default:
		http.Error(w, "", http.StatusBadRequest)
	}

}
