package client

import (
	"testing"
	"time"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/client/mocks"
	"github.com/JavaHutt/hashcash/internal/models"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestSolveChallenge(t *testing.T) {
	cfg := configs.Config{
		HashMaxIterations: 1e6,
		HashExpiration:    time.Hour,
	}
	client := &client{cfg: cfg}
	hashcash := models.Hashcash{
		Date: time.Now(),
		Bits: 1,
	}

	solved, err := client.solveChallenge(hashcash)
	assert.NoError(t, err)
	assert.NotNil(t, solved)
	assert.Greater(t, solved.Counter, 0)
}

func TestSolveChallengeExpired(t *testing.T) {
	cfg := configs.Config{
		HashMaxIterations: 1e6,
		HashExpiration:    -time.Hour,
	}
	client := &client{cfg: cfg}
	hashcash := models.Hashcash{
		Date: time.Now(),
		Bits: 1,
	}

	_, err := client.solveChallenge(hashcash)
	assert.ErrorIs(t, err, ErrDateCheck)
}

func TestSolveChallengeMaxIterExceed(t *testing.T) {
	cfg := configs.Config{
		HashMaxIterations: 1,
		HashExpiration:    time.Hour,
	}
	client := &client{cfg: cfg}
	hashcash := models.Hashcash{
		Date: time.Now(),
		Bits: 1,
	}

	_, err := client.solveChallenge(hashcash)
	assert.ErrorIs(t, err, ErrMaxIterExceed)
}

func TestRespondSolved(t *testing.T) {
	const ok = "OK"
	conn := &mocks.MockConn{}
	_, err := conn.ReadBuffer.WriteString(ok + string(models.EOFDelim))
	assert.NoError(t, err)

	client := &client{}
	hashcash := models.Hashcash{}

	resp, err := client.respondSolved(conn, hashcash)
	assert.NoError(t, err)
	assert.Equal(t, "OK", resp)
}

func TestRequestService(t *testing.T) {
	conn := &mocks.MockConn{}

	h := models.Hashcash{
		Ver:      1,
		Bits:     20,
		Date:     time.Now(),
		Resource: "127.0.0.1",
		Counter:  1,
	}
	conn.ReadBuffer.Write([]byte(h.String() + string(models.EOFDelim)))
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	client := &client{
		logger: logger.Sugar(),
	}

	hashcash, err := client.requestService(conn)
	assert.NoError(t, err)
	assert.NotNil(t, hashcash)

	assert.Equal(t, 20, hashcash.Bits)
	assert.Equal(t, 1, hashcash.Counter)

	request := conn.WriteBuffer.String()
	assert.Contains(t, request, models.MessageKindRequestChallenge)
}
