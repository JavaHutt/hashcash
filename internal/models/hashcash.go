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

type Hashcash struct {
	Ver      int
	Bits     int
	Date     time.Time
	Resource string
	Rand     string
	Counter  int
}

func (h Hashcash) String() string {
	encodedRand := base64.StdEncoding.EncodeToString([]byte(h.Rand))
	encodedCounter := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.Counter)))
	return fmt.Sprintf("%d:%d:%s:%s::%s:%s", h.Ver, h.Bits, h.Date.Format(dateLayout), h.Resource, encodedRand, encodedCounter)
}

func CheckHash(hash []byte, bits int) bool {
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
