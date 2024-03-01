package routes

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/jackc/pgx/v5"
)

// var (
// 	backends = []string{
// 		"http://localhost:8081",
// 		"http://localhost:8082",
// 	}
// 	rrIndex int
// )

type Server struct {
	addr  string
	conn  *pgx.Conn
	proxy *httputil.ReverseProxy
}

func NewServer(conn *pgx.Conn, apiID string) *Server {

	addr := fmt.Sprintf("http://api0%s:3000", apiID)
	url, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return &Server{
		addr:  addr,
		conn:  conn,
		proxy: proxy}
}

func (s *Server) GetAddress() string {
	return s.addr
}

func (s *Server) IsAlive() bool {
	return true
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}
