package main

type ServerPool struct {
	Servers []*Server
}

func NewServerPool() *ServerPool {
	return &ServerPool{
		Servers: make([]*Server, 0),
	}
}

func (sp *ServerPool) AddServer(server *Server) {
	sp.Servers = append(sp.Servers, server)
}

func (sp *ServerPool) GetAllServers() []*Server {
	return sp.Servers
}
