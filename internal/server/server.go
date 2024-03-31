package server

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/models"

	"go.uber.org/zap"
)

type server struct {
	cfg    configs.Config
	logger *zap.SugaredLogger
	wg     sync.WaitGroup
	li     net.Listener
}

// NewServer is a constructor
func NewServer(cfg configs.Config, logger *zap.SugaredLogger) *server {
	return &server{
		cfg:    cfg,
		logger: logger,
	}
}

// Listen starts listening for incoming TCP connections. Blocking operation
func (s *server) Listen(ctx context.Context) error {
	port := strconv.Itoa(s.cfg.Port)
	var err error
	if s.li, err = net.Listen("tcp", ":"+port); err != nil {
		return fmt.Errorf("error listening tcp connection: %w", err)
	}
	defer s.li.Close()

	s.logger.Infof("listening at %s", port)

	for {
		conn, err := s.li.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.logger.Info("server stopping as listener was closed")
				return nil
			}
			return fmt.Errorf("error accepting message: %w", err)
		}
		s.wg.Add(1)

		go func() {
			defer s.wg.Done()
			s.handleRequest(conn)
		}()
	}
}

// Shutdown gracefully shuts down the server.
func (s *server) Shutdown(ctx context.Context) error {
	if err := s.li.Close(); err != nil {
		s.logger.Errorf("failed to close listener: %v", err)
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		s.logger.Info("all connections closed")
		return nil
	}
}

func (s *server) handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
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
				return
			}
		case models.MessageKindSolvedChallenge:
			if err = s.verifySolved(conn, conn.RemoteAddr(), msg); err != nil {
				s.logger.Errorf("failed to verify solved challenge: %v", err)
			}
			return
		default:
			s.logger.Warnf("unknown message kind: %v", msg.Kind)
			return
		}
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

	if err = writeResp(conn, hashcash.String()); err != nil {
		return fmt.Errorf("failed to write hashcash string: %w", err)
	}

	return nil
}

func (s *server) verifySolved(w io.Writer, clientAddr net.Addr, msg models.Message) error {
	resource, _, err := net.SplitHostPort(clientAddr.String())
	if err != nil {
		return fmt.Errorf("failed to split host port: %w", err)
	}

	hashcash, err := models.ParseHashcash(msg.Hashcash)
	if err != nil {
		return fmt.Errorf("failed to parse hashcash: %w", err)
	}

	if hashcash.Resource != resource {
		return errors.New("hashcash resource and client addr aren't matching")
	}

	if !s.checkDate(hashcash.Date) {
		return errors.New("didn't pass date check")
	}

	hash := sha1.Sum([]byte(hashcash.String()))
	if !models.CheckHash(hash[:], hashcash.Bits) {
		return errors.New("hash does not meet the difficulty criteria")
	}

	if err := writeResp(w, getWidsom()); err != nil {
		return err
	}

	return nil
}

func (s *server) checkDate(date time.Time) bool {
	now := time.Now().UTC()
	return date.Before(now.Add(s.cfg.HashExpiration))
}

func writeResp(w io.Writer, response string) error {
	if _, err := w.Write([]byte(response + string(models.EOFDelim))); err != nil {
		return err
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
