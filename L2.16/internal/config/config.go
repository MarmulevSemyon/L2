package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

// Config хранит параметры запуска программы.
type Config struct {
	URL         string
	Depth       int
	OutputDir   string
	Concurrency int
}

// Parse разбирает аргументы командной строки и возвращает конфигурацию приложения.
func Parse(args []string) (Config, error) {
	var cfg Config

	fs := pflag.NewFlagSet("wget", pflag.ContinueOnError)

	fs.IntVarP(&cfg.Depth, "depth", "d", 1, "recursion depth")
	fs.StringVarP(&cfg.OutputDir, "output", "o", "output", "output directory")
	fs.IntVarP(&cfg.Concurrency, "concurrency", "c", 1, "number of concurrent downloads")

	if err := fs.Parse(args); err != nil {
		return Config{}, fmt.Errorf("parse flags: %w", err)
	}

	rest := fs.Args()
	if len(rest) == 0 {
		return Config{}, fmt.Errorf("не указан url")
	}
	if len(rest) > 1 {
		return Config{}, fmt.Errorf("слишком много аргументов: %v", rest[1:])
	}
	if cfg.Concurrency <= 0 {
		return Config{}, fmt.Errorf("concurrency must be > 0, got %d", cfg.Concurrency)
	}

	cfg.URL = rest[0]

	if cfg.Depth < 0 {
		return Config{}, fmt.Errorf("глубина должна быть >= 0, получено: %d", cfg.Depth)
	}

	return cfg, nil
}
