package main

import (
	"context"
	"log"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg, err := configs.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	ctx := context.Background()

	s := server.NewServer(*cfg, sugar)
	s.Listen(ctx)
}
