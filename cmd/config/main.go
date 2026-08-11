// Command config reads a TOML configuration file (config.toml) into a struct
// using github.com/pelletier/go-toml.
//
//	go run ./cmd/config
package main

import (
	"fmt"
	"log"
	"os"

	toml "github.com/pelletier/go-toml"
)

// Config mirrors the structure of config.toml.
type Config struct {
	Login struct {
		User     string
		Password string
	}
}

func main() {
	file, err := os.Open("config.toml")
	if err != nil {
		log.Fatalf("error: can't open config file - %s", err)
	}
	defer func() { _ = file.Close() }()

	cfg := &Config{}
	if err := toml.NewDecoder(file).Decode(cfg); err != nil {
		log.Fatalf("error: can't decode configuration file - %s", err)
	}

	fmt.Printf("%+v\n", cfg)
}
