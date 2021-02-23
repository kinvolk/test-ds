package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/kinvolk/test-ds/internal"
)

type Config struct {
	ClusterName      string `json:"cluster-name"`
	LocalPort        int    `json:"local-port"`
	NodeID           string `json:"node-id"`
	DiscoveryPort    int    `json:"discovery-port"`
	ControlPlaneName string `json:"control-plane-name"`
}

func NewConfig(cfgPath string) (Config, error) {
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, err
	}
	defer cfgFile.Close()

	cfg := Config{}
	if err := json.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	if cfg.ClusterName == "" {
		return errors.New("Empty cluster name")
	}
	if err := internal.ValidatePort(cfg.LocalPort, "local"); err != nil {
		return err
	}
	if cfg.NodeID == "" {
		return errors.New("Empty node ID")
	}
	if err := internal.ValidatePort(cfg.DiscoveryPort, "discovery"); err != nil {
		return err
	}
	if cfg.ControlPlaneName == "" {
		return errors.New("Empty control plane name")
	}
	return nil
}
