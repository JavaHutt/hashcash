package server

import (
	"net"
	"testing"

	"github.com/JavaHutt/hashcash/configs"
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
