package client

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/models"

	"go.uber.org/zap"
)

const eofDelim = '\n'

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
	wisdom, err := c.respondSolved(conn, *hashcash)
	if err != nil {
		return err
	}
	c.logger.Infof("Thou shalt hear the word of wisdom: %s 👼", wisdom)

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
		return nil, errors.New("didn't pass date check")
	}

	for range c.cfg.HashMaxIterations {
		hash := sha1.Sum([]byte(hashcash.String()))

		if checkHash(hash[:], hashcash.Bits) {
			return &hashcash, nil
		}

		hashcash.Counter++
	}

	return nil, errors.New("maximum iterations exceeded")
}

func (c *client) checkDate(date time.Time) bool {
	now := time.Now().UTC()
	return date.Before(now.Add(c.cfg.HashExpiration))
}

func checkHash(hash []byte, bits int) bool {
	zeroBytes := bits / 8
	zeroBits := bits % 8
	for i := range zeroBytes {
		if hash[i] != 0 {
			return false
		}
	}

	if zeroBits > 0 {
		mask := byte(0xFF << (8 - zeroBits))
		return hash[zeroBytes]&mask == 0
	}

	return true
}

func (c *client) respondSolved(conn net.Conn, hashcash models.Hashcash) (string, error) {
	msg := models.Message{
		Kind:     models.MessageKindSolvedChallenge,
		Hashcash: hashcash.String(),
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal solved message: %w", err)
	}

	if _, err = conn.Write(b); err != nil {
		return "", fmt.Errorf("failed to write solved message: %w", err)
	}

	resp, err := readResponse(conn)
	if err != nil {
		return "", fmt.Errorf("failed to read granted response: %w", err)
	}

	return string(resp), nil
}

func readResponse(r io.Reader) ([]byte, error) {
	reader := bufio.NewReader(r)
	resp, err := reader.ReadBytes(eofDelim)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes: %w", err)
	}
	resp = bytes.TrimSuffix(resp, []byte{eofDelim})

	return resp, nil
}

func getAddress(cfg configs.Config) string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
