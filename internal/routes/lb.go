package routes

import (
	"log"
	"net/http"
)

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []*Server
}

func NewLoadBalancer(port string, servers []*Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func (lb *LoadBalancer) rotate() *Server {

	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++

	return server
}

func (lb *LoadBalancer) GetPort() string {
	return lb.port
}

func (lb *LoadBalancer) GetServers() []*Server {
	return lb.servers
}

func (lb *LoadBalancer) ServeProxy(w http.ResponseWriter, r *http.Request) {
	s := lb.rotate()

	log.Print("API calling by the address: ", s.addr)
	s.Handler(w, r)
}
