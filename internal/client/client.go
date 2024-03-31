package client

import (
	"context"
	"fmt"

	"github.com/JavaHutt/hashcash/configs"
)

type client struct{}

func NewClient(cfg configs.Config) *client {
	return &client{}
}

func (c *client) Run(ctx context.Context) error {
	fmt.Println("client run")
	return nil
}
