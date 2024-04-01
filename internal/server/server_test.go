package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/models"
	"github.com/JavaHutt/hashcash/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestChooseChallenge(t *testing.T) {
	s := &server{
		cfg: configs.Config{
			HashBits:    20,
			HashCounter: 1,
		},
		logger: zap.NewNop().Sugar(),
	}

	mockConn := &mocks.MockConn{}
	clientAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	assert.NoError(t, err)

	err = s.chooseChallenge(mockConn, clientAddr)
	assert.NoError(t, err)

	writeBufferContents := mockConn.WriteBuffer.String()
	assert.NotEmpty(t, writeBufferContents)

	expectedResource := "127.0.0.1"
	assert.Contains(t, writeBufferContents, expectedResource)
}

func TestVerifySolvedHashExists(t *testing.T) {
	mockStore := mocks.NewMockStore()
	mockConn := &mocks.MockConn{}

	s := &server{
		store:  mockStore,
		logger: zap.NewNop().Sugar(),
	}

	ctx := context.Background()
	clientAddr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
	hashcashStr := "1:5:240401110844:127.0.0.1::NjAw:Mzg="
	msg := models.Message{Hashcash: hashcashStr}

	mockStore.Set(ctx, hashcashStr)

	err := s.verifySolved(ctx, mockConn, clientAddr, msg)
	assert.ErrorIs(t, err, ErrHashcashExists)
}

func TestVerifySolvedAddrMismatch(t *testing.T) {
	mockStore := mocks.NewMockStore()
	mockConn := &mocks.MockConn{}

	s := &server{
		store:  mockStore,
		logger: zap.NewNop().Sugar(),
	}

	ctx := context.Background()
	clientAddr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
	hashcashStr := "1:5:240401110844:..1::NjAw:Mzg="
	msg := models.Message{Hashcash: hashcashStr}

	err := s.verifySolved(ctx, mockConn, clientAddr, msg)
	assert.ErrorIs(t, err, ErrAddrMismatch)
}

func TestVerifySolvedDateCheck(t *testing.T) {
	mockStore := mocks.NewMockStore()
	mockConn := &mocks.MockConn{}

	s := &server{
		cfg: configs.Config{
			HashExpiration: time.Hour,
		},
		store:  mockStore,
		logger: zap.NewNop().Sugar(),
	}

	ctx := context.Background()
	clientAddr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
	hashcashStr := "1:5:200401110844:127.0.0.1::NjAw:Mzg="
	msg := models.Message{Hashcash: hashcashStr}

	err := s.verifySolved(ctx, mockConn, clientAddr, msg)
	assert.ErrorIs(t, err, ErrDateCheck)
}

func Test_normalizeIPAddress(t *testing.T) {
	tests := []struct {
		name string
		addr string
		want string
	}{
		{
			name: "ipv4",
			addr: "192.168.1.1:3000",
			want: "192.168.1.1",
		},
		{
			name: "ipv6",
			addr: "[::1]:51178",
			want: "..1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeIPAddress(tt.addr); got != tt.want {
				t.Errorf("normalizeIPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
