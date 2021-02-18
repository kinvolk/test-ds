package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := progMain(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func progMain() error {
	cfgPath := ""
	flag.StringVar(&cfgPath, "config", "", "path to json config")
	flag.Parse()

	cfg, err := NewConfig(cfgPath)
	if err != nil {
		return err
	}
	ctx := context.Background()
	server := NewServer(ctx, cfg.Hash)
	return server.Run(cfg.Port)
}
