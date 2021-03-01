package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/kinvolk/test-ds/internal"
)

type Config struct {
	Hash string `json:"hash"`
	Port int    `json:"port"`
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
	if cfg.Hash == "" {
		return errors.New("empty hash")
	}
	m := map[string]struct{}{}
	m["sha256"] = struct{}{}
	m["md5"] = struct{}{}
	if _, ok := m[cfg.Hash]; !ok {
		return fmt.Errorf("invalid hash %s", cfg.Hash)
	}
	if err := internal.ValidatePort(cfg.Port, "listen"); err != nil {
		return err
	}
	return nil
}
