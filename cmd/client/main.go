package main

import (
	"context"
	"log"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/client"
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

	c := client.NewClient(*cfg, sugar)
	if err = c.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
