package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type Storage struct {
	rootDir string
}

func New(rootDir string) *Storage {
	return &Storage{
		rootDir: rootDir,
	}
}

func (s *Storage) RootDir() string {
	return s.rootDir
}

func (s *Storage) SaveFile(localPath string, data []byte) error {
	dir := filepath.Dir(localPath)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directories: %w", err)
	}

	if err := os.WriteFile(localPath, data, 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
