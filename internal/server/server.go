package server

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/models"

	"go.uber.org/zap"
)

type store interface {
	Set(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type server struct {
	cfg      configs.Config
	logger   *zap.SugaredLogger
	wg       sync.WaitGroup
	li       net.Listener
	store    store
	rTimeout time.Duration
	wTimeout time.Duration
}

// NewServer is a constructor
func NewServer(cfg configs.Config, logger *zap.SugaredLogger, store store) *server {
	s := server{
		cfg:    cfg,
		logger: logger,
		store:  store,
	}

	if cfg.ReadTimeout != 0 {
		s.rTimeout = cfg.ReadTimeout
	}

	if cfg.WriteTimeout != 0 {
		s.wTimeout = cfg.WriteTimeout
	}

	return &s
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
			s.handleRequest(ctx, conn)
		}()
	}
}

// Shutdown gracefully shuts down the server
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

func (s *server) handleRequest(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	for {
		msg, err := s.decodeMessage(conn)
		if err != nil {
			s.logger.Errorf("error decoding message: %v", err)
			return
		}

		s.logger.Infof("received message: %v", msg)

		switch msg.Kind {
		case models.MessageKindRequestChallenge:
			if err = s.chooseChallenge(conn); err != nil {
				s.logger.Errorf("failed to choose challenge: %v", err)
				return
			}
		case models.MessageKindSolvedChallenge:
			if err = s.verifySolved(ctx, conn, *msg); err != nil {
				s.logger.Errorf("failed to verify solved challenge: %v", err)
			}
			return
		default:
			s.logger.Warnf("unknown message kind: %v", msg.Kind)
			return
		}
	}
}

func (s *server) chooseChallenge(conn net.Conn) error {
	resource := normalizeIPAddress(conn.RemoteAddr().String())

	hashcash := models.Hashcash{
		Ver:      1,
		Bits:     s.cfg.HashBits,
		Date:     time.Now().UTC(),
		Resource: resource,
		Rand:     fmt.Sprintf("%d", rand.Intn(1e3)),
		Counter:  s.cfg.HashCounter,
	}

	if err := s.writeResp(conn, hashcash.String()); err != nil {
		return fmt.Errorf("failed to write hashcash string: %w", err)
	}

	return nil
}

func (s *server) verifySolved(ctx context.Context, conn net.Conn, msg models.Message) error {
	resource := normalizeIPAddress(conn.RemoteAddr().String())

	exists, err := s.store.Exists(ctx, msg.Hashcash)
	if err != nil {
		return fmt.Errorf("failed to get hashcash from the store: %w", err)
	}
	if exists {
		return ErrHashcashExists
	}

	hashcash, err := models.ParseHashcash(msg.Hashcash)
	if err != nil {
		return fmt.Errorf("failed to parse hashcash: %w", err)
	}

	if hashcash.Resource != resource {
		return ErrAddrMismatch
	}

	if !s.checkDate(hashcash.Date) {
		return ErrDateCheck
	}

	hash := sha1.Sum([]byte(hashcash.String()))
	if !models.CheckHash(hash[:], hashcash.Bits) {
		return ErrCheckHash
	}

	s.logger.Infof("hash %s has passed all checks!", msg.Hashcash)

	if err = s.store.Set(ctx, msg.Hashcash); err != nil {
		return fmt.Errorf("failed to save hash in the storage: %w", err)
	}

	return s.writeResp(conn, getRandomWisdom())
}

func (s *server) checkDate(date time.Time) bool {
	now := time.Now().UTC()
	return now.Before(date.Add(s.cfg.HashExpiration))
}

func (s *server) writeResp(conn net.Conn, response string) error {
	err := conn.SetWriteDeadline(time.Now().Add(s.wTimeout))
	if err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	if _, err = conn.Write([]byte(response + string(models.EOFDelim))); err != nil {
		return fmt.Errorf("failed to write to connection: %w", err)
	}

	if err := conn.SetWriteDeadline(time.Time{}); err != nil {
		return fmt.Errorf("failed to reset write deadline: %w", err)
	}

	return nil
}

func (s *server) decodeMessage(conn net.Conn) (*models.Message, error) {
	if err := conn.SetReadDeadline(time.Now().Add(s.rTimeout)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	var msg models.Message
	err := json.NewDecoder(conn).Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("error decoding message: %w", err)
	}

	if err = conn.SetReadDeadline(time.Time{}); err != nil {
		return nil, fmt.Errorf("failed to reset read deadline: %w", err)
	}

	return &msg, nil
}

func normalizeIPAddress(addr string) string {
	host, _, splitErr := net.SplitHostPort(addr)
	if splitErr == nil {
		addr = host
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		return addr
	}

	if ip.To4() == nil {
		// this is my little trick to create hashcash string without confusing `:` symbols
		return strings.ReplaceAll(ip.String(), ":", ".")
	}

	return ip.String()
}
