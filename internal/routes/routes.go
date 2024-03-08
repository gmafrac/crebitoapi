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

	client := &db.Client{ID: id}

	switch r.Method {

	case http.MethodGet:
		if pathParts[2] != "extrato" {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		s.Extract(w, client)

	case http.MethodPost:
		if pathParts[2] != "transacoes" {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		s.Transaction(w, r, client)

	default:
		http.Error(w, "", http.StatusBadRequest)
	}

}
