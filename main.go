package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcplog "github.com/envoyproxy/go-control-plane/pkg/log"
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

	logger := log.New(os.Stderr, fmt.Sprintf("|%s| ", cfg.ControlPlaneName), log.Lmsgprefix)
	cacheLogger := gcplog.LoggerFuncs{
		DebugFunc: makeLogFunc(logger, "DEBUG"),
		InfoFunc:  makeLogFunc(logger, "INFO"),
		WarnFunc:  makeLogFunc(logger, "WARN"),
		ErrorFunc: makeLogFunc(logger, "ERROR"),
	}
	scache := cache.NewSnapshotCache(false, cache.IDHash{}, cacheLogger)
	if err := scache.SetSnapshot(cfg.NodeID, snapshot); err != nil {
		return err
	}

	ctx := context.Background()
	server := NewServer(ctx, scache, logger, cfg.ControlPlaneName)
	return server.Run(ctx, cfg.DiscoveryPort)
}

func makeLogFunc(logger *log.Logger, level string) func(string, ...interface{}) {
	return func(format string, args ...interface{}) {
		logFunc(logger, level, format, args...)
	}
}

func logFunc(logger *log.Logger, level, format string, args ...interface{}) {
	logger.Printf("%s: "+format, append([]interface{}{level}, args...)...)
}
