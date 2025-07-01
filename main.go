package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log/slog"
	"net"
	"sync"
)

const (
	apiKeyProduce int16 = 0
	apiKeyFetch   int16 = 1
)

type Partition struct {
	log           [][]byte
	highWatermark int64
	//..
}
type Server struct {
	mu     sync.Mutex
	topics map[string]map[int32]*Partition
	ln     net.Listener
}

func NewServer() *Server {
	return &Server{
		topics: make(map[string]map[int32]*Partition),
	}
}

func (s *Server) getOrCreatePartition(topic string, partitionID int32) *Partition {
	if _, ok := s.topics[topic]; !ok {
		s.topics[topic] = make(map[int32]*Partition)
	}
	if _, ok := s.topics[topic][partitionID]; !ok {
		s.topics[topic][partitionID] = &Partition{log: make([][]byte, 0)}
	}
	return s.topics[topic][partitionID]
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
	slog.Info("New Connection", "remoteAddr", conn.RemoteAddr())

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

		//bytes.Reader allows for easier parsing: https://pkg.go.dev/bytes#NewReader
		r := bytes.NewReader(requestBytes)
		//parsing header using bytReader
		apiKey, _ := readInt16(r)
		apiVersion, _ := readInt16(r) //routing for later
		correlationID, _ := readInt32(r)
		// clientID, _ := readString(r)

		clientIDLen, _ := readInt16(r)
		r.Seek(int64(clientIDLen), io.SeekCurrent)
		//clientID := string(requestBytes[10:headerEndOffset])

		slog.Info("Received reqest", "apiKey", apiKey, "apiVersion", apiVersion, "correlationID", correlationID)

		switch apiKey {
		case apiKeyProduce:
			s.handleProduce(conn, correlationID, r)
		case apiKeyFetch:
			s.handleFetch(conn, correlationID, r)
		default:
			slog.Warn("Received unknown API key", "apiKey", apiKey)
		}
	}
}

// parses produce request, adds message to buffer, sends response
func (s *Server) handleProduce(conn net.Conn, correlationID int32, r *bytes.Reader) {
	return
}

// parses fetch request, reads from buffer, sends messages back to client
func (s *Server) handleFetch(conn net.Conn, correlationID int32, r *bytes.Reader) {
	return
}

func main() {
	server := NewServer()
	server.Listen()
}
