package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/log"
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

	snapshot, err := GetSnapshot(cfg.ClusterName, cfg.LocalPort)
	if err != nil {
		return err
	}

	logger := log.LoggerFuncs{
		DebugFunc: makeLogFunc("DEBUG"),
		InfoFunc:  makeLogFunc("INFO"),
		WarnFunc:  makeLogFunc("WARN"),
		ErrorFunc: makeLogFunc("ERROR"),
	}
	scache := cache.NewSnapshotCache(false, cache.IDHash{}, logger)
	if err := scache.SetSnapshot(cfg.NodeID, snapshot); err != nil {
		return err
	}

	ctx := context.Background()
	server := NewServer(ctx, scache)
	return server.Run(ctx, cfg.DiscoveryPort)
}

func makeLogFunc(level string) func(string, ...interface{}) {
	return func(format string, args ...interface{}) {
		logFunc(level, format, args...)
	}
}

func logFunc(level, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: "+format+"\n", append([]interface{}{level}, args...))
}
