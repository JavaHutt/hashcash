package main

import (
	"context"
	"log"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/JavaHutt/hashcash/internal/client"
)

func main() {
	cfg, err := configs.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	c := client.NewClient(*cfg)
	if err = c.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
