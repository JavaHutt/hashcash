package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/models"
)

type client struct {
	cfg configs.Config
}

func NewClient(cfg configs.Config) *client {
	return &client{
		cfg: cfg,
	}
}

func (c *client) Run(_ context.Context) error {
	conn, err := net.Dial("tcp", getAddress(c.cfg))
	if err != nil {
		return fmt.Errorf("couldn't connect to server: %w", err)
	}

	defer conn.Close()

	// 1. Request service
	hashcash, err := c.requestService(conn)
	if err != nil {
		return err
	}

	fmt.Printf("hashcash: %+v", *hashcash)

	// 2. SolveChallenge

	// 3. ResponseSolved
	return nil
}

func (c *client) requestService(conn net.Conn) (*models.Hashcash, error) {
	msg := models.Message{Kind: models.MessageKindRequestChallenge}
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request message: %w", err)
	}

	if _, err = conn.Write(b); err != nil {
		return nil, fmt.Errorf("failed to write request message: %w", err)
	}

	resp, err := io.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}

	hashcash, err := models.ParseHashcash(string(resp))
	if err != nil {
		return nil, fmt.Errorf("failed to parse hashcash: %w", err)
	}

	return hashcash, nil
}

func getAddress(cfg configs.Config) string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
