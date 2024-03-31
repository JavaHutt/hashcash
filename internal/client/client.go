package client

import (
	"context"
	"encoding/json"
	"fmt"
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

func (c *client) Run(ctx context.Context) error {
	conn, err := net.Dial("tcp", getAddress(c.cfg))
	if err != nil {
		return fmt.Errorf("couldn't connect to server: %w", err)
	}

	defer conn.Close()

	// 1. Request service
	if err = c.requestService(ctx, conn); err != nil {
		return err
	}

	// 2. SolveChallenge

	// 3. ResponseSolved
	return nil
}

func (c *client) requestService(ctx context.Context, conn net.Conn) error {
	msg := models.Message{Kind: models.MessageKindRequestChallenge}
	b, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal request message: %w", err)
	}

	if _, err = conn.Write(b); err != nil {
		return fmt.Errorf("failed to write request message: %w", err)
	}

	return nil
}

func getAddress(cfg configs.Config) string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
