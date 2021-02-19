package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	if len(flag.Args()) != 2 {
		return errors.New("expected two parameters, for the server name and message to send")
	}
	serverName := flag.Arg(0)
	if serverName == "" {
		return errors.New("the server parameter can not be empty")
	}
	message := flag.Arg(1)
	if message == "" {
		return errors.New("message can not be empty")
	}

	cfg, err := NewConfig(cfgPath)
	if err != nil {
		return err
	}

	port, ok := cfg.ServerPorts[serverName]
	if !ok {
		return fmt.Errorf("unknown server %s", serverName)
	}

	address := "http://" + net.JoinHostPort("localhost", strconv.Itoa(port))
	messageReader := strings.NewReader(message)
	response, err := http.Post(address, "application/octet-stream", messageReader)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	_, _ = io.Copy(os.Stdout, response.Body)
	return nil
}
