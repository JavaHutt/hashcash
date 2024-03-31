package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/models"
)

type server struct {
	cfg configs.Config
}

func NewServer(cfg configs.Config) *server {
	return &server{
		cfg: cfg,
	}
}

func (s *server) Listen(ctx context.Context) error {
	port := strconv.Itoa(s.cfg.Port)
	li, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("error listening tcp connection: %w", err)
	}
	defer li.Close()

	fmt.Printf("Listening on %s", port)
	for {
		conn, err := li.Accept()
		if err != nil {
			return fmt.Errorf("error accepting message: %w", err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	var msg models.Message
	err := json.NewDecoder(conn).Decode(&msg)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	fmt.Println("Received:", msg)
	if _, err = conn.Write([]byte("Message received.\n")); err != nil {
		log.Printf("failed to write to connection: %w", err)
	}
}
