package config

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

// Config хранит параметры запуска программы.
type Config struct {
	Timeout time.Duration
	Host    string
	Port    string
}

// Parse разбирает аргументы командной строки и возвращает конфигурацию приложения.
func Parse(args []string) (Config, error) {
	var cfg Config

	fs := pflag.NewFlagSet("telnet", pflag.ContinueOnError)

	fs.DurationVarP(&cfg.Timeout, "timeout", "T", 10*time.Second, "connection timeout")

	if err := fs.Parse(args); err != nil {
		return Config{}, fmt.Errorf("parse flags: %w", err)
	}
	rest := fs.Args()
	if len(rest) != 2 {
		return Config{}, fmt.Errorf("usage: telnet [--timeout=10s] host port")
	}
	cfg.Host = rest[0]
	cfg.Port = rest[1]
	return cfg, nil
}
