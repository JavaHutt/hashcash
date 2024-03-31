package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/models"

	"go.uber.org/zap"
)

type server struct {
	cfg    configs.Config
	logger *zap.SugaredLogger
}

func NewServer(cfg configs.Config, logger *zap.SugaredLogger) *server {
	return &server{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *server) Listen(ctx context.Context) error {
	port := strconv.Itoa(s.cfg.Port)
	li, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("error listening tcp connection: %w", err)
	}
	defer li.Close()

	s.logger.Infof("listening at %s", port)

	for {
		conn, err := li.Accept()
		if err != nil {
			return fmt.Errorf("error accepting message: %w", err)
		}

		go s.handleRequest(conn)
	}
}

func (s *server) handleRequest(conn net.Conn) {
	defer conn.Close()

	msg, err := decodeMessage(conn)
	if err != nil {
		s.logger.Errorf("error decoding message: %v", err)
		return
	}

	s.logger.Infof("received message: %v", msg)

	switch msg.Kind {
	case models.MessageKindRequestChallenge:
		if err = s.chooseChallenge(conn, conn.RemoteAddr()); err != nil {
			s.logger.Errorf("failed to choose challenge: %v", err)
		}
	default:
		s.logger.Warnf("unknown message kind: %v", msg.Kind)
	}
}

func (s *server) chooseChallenge(conn io.Writer, clientAddr net.Addr) error {
	resource, _, err := net.SplitHostPort(clientAddr.String())
	if err != nil {
		return fmt.Errorf("failed to split host port: %w", err)
	}

	hashcash := models.Hashcash{
		Ver:      1,
		Bits:     s.cfg.HashBits,
		Date:     time.Now().UTC(),
		Resource: resource,
		Rand:     fmt.Sprintf("%d", rand.Intn(1e3)),
		Counter:  3,
	}

	if _, err = conn.Write([]byte(hashcash.String())); err != nil {
		return fmt.Errorf("failed to write hashcash string: %w", err)
	}

	return nil
}

func decodeMessage(conn net.Conn) (models.Message, error) {
	var msg models.Message
	err := json.NewDecoder(conn).Decode(&msg)
	if err != nil {
		return models.Message{}, fmt.Errorf("error decoding message: %w", err)
	}

	return msg, nil
}
