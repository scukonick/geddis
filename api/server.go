package api

import (
	"bufio"
	"log"
	"net"
	"strconv"

	"github.com/scukonick/geddis/db"
)

// server serves
type server struct {
	l      net.Listener
	store  *db.GeddisStore
	logger *log.Logger
}

func newServer(l net.Listener, store *db.GeddisStore, logger *log.Logger) *server {
	return &server{
		l:      l,
		store:  store,
		logger: logger,
	}
}

func (s *server) serve() {
	for {
		// Wait for a connection.
		conn, err := s.l.Accept()
		if err != nil {
			s.logger.Printf("ERR: failed to accept connection: %v", err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go s.handleConnect(conn)
	}
}

func (s *server) handleConnect(c net.Conn) {
	sAddr := c.RemoteAddr().String()

	s.logger.Printf("INFO: new connect from %s", sAddr)

	scanner := bufio.NewScanner(c)

	for scanner.Scan() {
		line := scanner.Text()
		line, err := bufReader.ReadLi('\n')
		if err != nil {
			s.logger.Printf("ERR: read from client %s failed: %v",
				sAddr, err)
			err = c.Close()
			if err != nil {
				s.logger.Printf("ERR: failed to close conn with client %s: %v", sAddr, err)
			}
			return
		}

		unquated, err := strconv.Unquote(line)
		if err == strconv.ErrSyntax {
			// line is not quoted
		}
	}

}
