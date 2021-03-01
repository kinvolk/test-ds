package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/kinvolk/test-ds/internal"
)

type Config struct {
	ServerPorts  map[string]int `json:"server-ports"`
	HostOverride string `json:"host-override"`
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
	if len(cfg.ServerPorts) == 0 {
		return errors.New("no server ports specified")
	}
	for name, port := range cfg.ServerPorts {
		if name == "" {
			return errors.New("empty server name")
		}
		if err := internal.ValidatePort(port, name); err != nil {
			return err
		}
	}
	return nil
}
