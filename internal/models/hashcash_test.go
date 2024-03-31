package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHashcash(t *testing.T) {
	dt, err := time.Parse(time.DateOnly, "2011-06-12")
	assert.NoError(t, err)

	h := Hashcash{
		Ver:      1,
		Bits:     5,
		Date:     dt,
		Resource: "127.0.0.1",
		Rand:     "random",
		Counter:  5,
	}

	hashString := h.String()
	want := "1:5:110612000000:127.0.0.1::cmFuZG9t:NQ=="
	assert.Equal(t, want, hashString)

	parsedHash, err := ParseHashcash(hashString)
	assert.NoError(t, err)
	assert.Equal(t, h, *parsedHash)
}
