package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go captureSigs(sigs, sugar, cancel)

	c := client.NewClient(*cfg, sugar)
	for {
		if err = c.Run(ctx); err != nil {
			log.Fatal(err)
		}

		select {
		case <-time.After(time.Second * 5):
			continue
		case <-ctx.Done():
			sugar.Info("context canceled, shutting down gracefully...")
			return
		}
	}
}

func captureSigs(sigs chan os.Signal, logger *zap.SugaredLogger, cancel context.CancelFunc) {
	sig := <-sigs
	logger.Infof("Received signal: %s", sig)
	cancel()
}
