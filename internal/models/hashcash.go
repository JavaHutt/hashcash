package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "060102150405"

// Hashcash is a object representation of hashcash string
type Hashcash struct {
	Ver      int
	Bits     int
	Date     time.Time
	Resource string
	Rand     string
	Counter  int
}

// String encodes hashcash struct into string
func (h Hashcash) String() string {
	encodedRand := base64.StdEncoding.EncodeToString([]byte(h.Rand))
	encodedCounter := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.Counter)))
	return fmt.Sprintf("%d:%d:%s:%s::%s:%s", h.Ver, h.Bits, h.Date.Format(dateLayout), h.Resource, encodedRand, encodedCounter)
}

// CheckHash verifies if the given hash satisfies the specified difficulty level.
// The difficulty is defined by the number of leading zero bits in the hash.
// The hash is considered to satisfy the difficulty if it has at least `bits` number of leading zero bits.
// Parameters:
// - hash: The hash to check, in bytes.
// - bits: The difficulty level, specified as the number of leading zero bits required.
// Returns true if the hash satisfies the difficulty level, false otherwise.
func CheckHash(hash []byte, bits int) bool {
	zeroBytes := bits / 8 // Calculate the number of fully zero bytes required for the given difficulty.
	zeroBits := bits % 8  // Calculate the remaining zero bits required in the next byte.
	for i := range zeroBytes {
		if hash[i] != 0 {
			return false
		}
	}

	if zeroBits > 0 {
		mask := byte(0xFF << (8 - zeroBits)) // Create a mask for the remaining zero bits.
		return hash[zeroBytes]&mask == 0     // Check if the remaining bits are zero using the mask.
	}

	return true
}

// ParseHashcash parses string to a hashcash struct
func ParseHashcash(h string) (*Hashcash, error) {
	parts := strings.SplitN(h, "::", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid hashcash format")
	}

	initialParts := strings.Split(parts[0], ":")
	if len(initialParts) != 4 {
		return nil, errors.New("invalid initial hashcash segment")
	}

	ver, err := strconv.Atoi(initialParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid version: %w", err)
	}

	bits, err := strconv.Atoi(initialParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid bits: %w", err)
	}

	date, err := time.Parse(dateLayout, initialParts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}

	resource := initialParts[3]

	finalParts := strings.Split(parts[1], ":")
	if len(finalParts) != 2 {
		return nil, errors.New("invalid final hashcash segment")
	}

	randBytes, err := base64.StdEncoding.DecodeString(finalParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid rand: %w", err)
	}
	rand := string(randBytes)

	counterB, err := base64.StdEncoding.DecodeString(finalParts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode counter: %w", err)
	}
	counter, err := strconv.Atoi(string(counterB))
	if err != nil {
		return nil, fmt.Errorf("invalid counter: %w", err)
	}

	return &Hashcash{
		Ver:      ver,
		Bits:     bits,
		Date:     date,
		Resource: resource,
		Rand:     rand,
		Counter:  counter,
	}, nil
}
