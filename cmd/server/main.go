package main

import (
	"context"
	"log"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/server"
)

func main() {
	cfg, err := configs.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	s := server.NewServer(*cfg)
	s.Run(ctx)
}
