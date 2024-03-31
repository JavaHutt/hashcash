package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := server.NewServer(*cfg, sugar)
	go func() {
		if err = s.Listen(ctx); err != nil {
			sugar.Fatalf("failed to start listening: %v", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	captureSigs(sigs, sugar, cancel)

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := s.Shutdown(ctxShutdown); err != nil {
		sugar.Fatalf("failed to shutdown server: %v", err)
	}

	sugar.Info("server gracefully stopped")
}

func captureSigs(sigs chan os.Signal, logger *zap.SugaredLogger, cancel context.CancelFunc) {
	sig := <-sigs
	logger.Infof("Received signal: %s", sig)
	cancel()
}
