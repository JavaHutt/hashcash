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

func TestCheckHash(t *testing.T) {
	tests := []struct {
		name  string
		hash  []byte
		bits  int
		valid bool
	}{
		{
			name:  "Zero hash, 0 bits",
			hash:  []byte{0, 0, 0, 0},
			bits:  0,
			valid: true,
		},
		{
			name:  "Zero hash, 8 bits",
			hash:  []byte{0, 0, 0, 0},
			bits:  8,
			valid: true,
		},
		{
			name:  "Zero hash, 9 bits",
			hash:  []byte{0, 0, 0, 0},
			bits:  9,
			valid: true,
		},
		{
			name:  "Non-zero hash, 8 bits (invalid)",
			hash:  []byte{1, 0, 0, 0},
			bits:  8,
			valid: false,
		},
		{
			name:  "Non-zero hash, 9 bits (valid)",
			hash:  []byte{0, 1, 0, 0}, // Second byte does not matter for 9 bits.
			bits:  9,
			valid: true,
		},
		{
			name:  "Non-zero hash, 16 bits (invalid)",
			hash:  []byte{0, 1, 0, 0},
			bits:  16,
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := CheckHash(test.hash, test.bits); got != test.valid {
				t.Errorf("CheckHash() = %v, want %v", got, test.valid)
			}
		})
	}
}
