package main

import (
	"encoding/binary"
	"io"
	"log/slog"
	"net"
	"sync"
)

const (
	apiKeyProduce int16 = 0
	apiKeyFetch   int16 = 0
)

type Message struct {
	data []byte
	//..
}
type Server struct {
	mu          sync.Mutex
	consOffsets map[string]int
	buffer      []Message
	ln          net.Listener
}

func NewServer() *Server {
	return &Server{
		consOffsets: make(map[string]int),
		buffer:      make([]Message, 0, 1024),
	}
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
	slog.Info("Listening on :9092")

	for {
		connection, err := ln.Accept()
		if err != nil {
			if err == io.EOF {
				slog.Error("Error accepting connection", "err", err)
				continue
			}
		}
		go s.handleConn(connection)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	slog.Info("New Connection at ", conn.RemoteAddr())

	for {
		//1. Read the Request (Size: 4 bytes)
		sizeBytes := make([]byte, 4)
		if _, err := io.ReadFull(conn, sizeBytes); err != nil {
			if err == io.EOF {
				slog.Info("Client closed connection", "remoteAddr", conn.RemoteAddr())
			} else {
				slog.Error("Could not read request", "err", err)
			}
			return
		}
		requestSize := binary.BigEndian.Uint32(sizeBytes)

		requestBytes := make([]byte, requestSize)
		if _, err := io.ReadFull(conn, requestBytes); err != nil {
			slog.Error("could not read full request", "err", err)
			return
		}

		//apiKey := int16(binary.BigEndian.Uint16(requestBytes[0:2]))
		//apiVersion := int16(binary.BigEndian.Uint16(requestBytes[2:4]))
		//correlationID := int32(binary.BigEndian.Uint32(requestBytes[4:8]))

	}
}

func main() {
	server := NewServer()
	server.Listen()
}
