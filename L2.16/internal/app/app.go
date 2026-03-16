package app

import (
	"fmt"
	"time"

	"wget/internal/config"
	"wget/internal/crawler"
	"wget/internal/fetcher"
	"wget/internal/storage"
)

func Run(cfg config.Config) error {
	fmt.Printf("start wget\n")
	fmt.Printf("url: %s\n", cfg.URL)
	fmt.Printf("depth: %d\n", cfg.Depth)
	fmt.Printf("output: %s\n", cfg.OutputDir)

	f := fetcher.New(10 * time.Second)
	s := storage.New(cfg.OutputDir)

	c, err := crawler.New(f, s, cfg.URL)
	if err != nil {
		return fmt.Errorf("create crawler: %w", err)
	}

	if err := c.Run(cfg.URL, cfg.Depth, cfg.Concurrency); err != nil {
		return fmt.Errorf("run crawler: %w", err)
	}
	if c.EntryPath() != "" {
		fmt.Printf("open this file offline: %s\n", c.EntryPath())
	}
	return nil
}
