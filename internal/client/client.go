package client

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/models"

	"go.uber.org/zap"
)

type client struct {
	cfg    configs.Config
	logger *zap.SugaredLogger
}

func NewClient(cfg configs.Config, logger *zap.SugaredLogger) *client {
	return &client{
		cfg:    cfg,
		logger: logger,
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

	// 2. Solve challenge
	solved, err := c.solveChallenge(*hashcash)
	if err != nil {
		return err
	}
	c.logger.Infof("counter before: %d, after solved: %d", hashcash.Counter, solved.Counter)

	// 3. Respond solved
	wisdom, err := c.respondSolved(conn, *solved)
	if err != nil {
		return err
	}
	c.logger.Infof("Thou shalt hear the word of wisdom: %s 👼", wisdom)

	return nil
}

func (c *client) requestService(conn io.ReadWriter) (*models.Hashcash, error) {
	msg := models.Message{Kind: models.MessageKindRequestChallenge}

	if err := writeJSONResp(conn, msg); err != nil {
		return nil, fmt.Errorf("failed to write request message: %w", err)
	}

	resp, err := readResponse(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read challenge response: %w", err)
	}

	hashcash, err := models.ParseHashcash(string(resp))
	if err != nil {
		return nil, fmt.Errorf("failed to parse hashcash: %w", err)
	}

	return hashcash, nil
}

func (c *client) solveChallenge(hashcash models.Hashcash) (*models.Hashcash, error) {
	if !c.checkDate(hashcash.Date) {
		return nil, ErrDateCheck
	}

	for range c.cfg.HashMaxIterations {
		hash := sha1.Sum([]byte(hashcash.String()))

		if models.CheckHash(hash[:], hashcash.Bits) {
			return &hashcash, nil
		}

		hashcash.Counter++
	}

	return nil, ErrMaxIterExceed
}

func (c *client) checkDate(date time.Time) bool {
	now := time.Now().UTC()
	expiration := date.Add(c.cfg.HashExpiration)
	return now.Before(expiration) || now.Equal(expiration)
}

func (c *client) respondSolved(conn io.ReadWriter, hashcash models.Hashcash) (string, error) {
	msg := models.Message{
		Kind:     models.MessageKindSolvedChallenge,
		Hashcash: hashcash.String(),
	}

	if err := writeJSONResp(conn, msg); err != nil {
		return "", fmt.Errorf("failed to write solved message: %w", err)
	}

	resp, err := readResponse(conn)
	if err != nil {
		return "", fmt.Errorf("failed to read granted response: %w", err)
	}

	return string(resp), nil
}

func writeJSONResp(w io.Writer, response any) error {
	b, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	if _, err = w.Write(b); err != nil {
		return fmt.Errorf("failed to write response message: %w", err)
	}

	return nil
}

func readResponse(r io.Reader) ([]byte, error) {
	reader := bufio.NewReader(r)
	resp, err := reader.ReadBytes(models.EOFDelim)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes: %w", err)
	}
	resp = bytes.TrimSuffix(resp, []byte{models.EOFDelim})

	return resp, nil
}

func getAddress(cfg configs.Config) string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
