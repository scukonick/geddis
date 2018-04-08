package api

import "github.com/scukonick/geddis/db"

// TCPServer implements TCP API for GeddisStorage
type TCPServer struct {
	store *db.GeddisStore
}

// NewTCPServer returns newly initialized *TCPServer
func NewTCPServer(store *db.GeddisStore) *TCPServer {
	return &TCPServer{
		store: store,
	}
}

func (s *TCPServer) Listen(addr string) {

}
