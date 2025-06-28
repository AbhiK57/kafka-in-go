package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
)

type Server struct {
	ln net.Listener
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	return nil
}

func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", ":9092")
	if err != nil {
		return err
	}
	s.ln = ln
	for {
		connection, err := ln.Accept()
		if err != nil {
			if err == io.EOF {
				return err
			}
			slog.Error("Error with server accepting", "err", err)
		}
		go s.handleConn(connection)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	fmt.Println("New Connection at ", conn.RemoteAddr())

}

func main() {
	server := NewServer()
	server.Listen()
}
