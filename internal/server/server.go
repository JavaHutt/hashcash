package server

import (
	"context"
	"fmt"

	"github.com/JavaHutt/hashcash/configs"
)

type server struct{}

func NewServer(cfg configs.Config) *server {
	return &server{}
}

func (s *server) Run(ctx context.Context) error {
	fmt.Println("server run")
	return nil
}
