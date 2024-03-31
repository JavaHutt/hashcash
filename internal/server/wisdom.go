package server

import (
	_ "embed"
	"math/rand"
	"strings"
	"time"
)

//go:embed data/wisdom.txt
var wisdom string

func getWidsom() string {
	quotes := strings.Split(wisdom, "\n")

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randomIndex := r.Intn(len(quotes))
	return quotes[randomIndex]
}
