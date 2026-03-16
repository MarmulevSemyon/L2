package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// Storage отвечает за сохранение файлов в локальную файловую систему.
type Storage struct {
	rootDir string
}

// Storage отвечает за сохранение файлов в локальную файловую систему.
func New(rootDir string) *Storage {
	return &Storage{
		rootDir: rootDir,
	}
}

// New создаёт новый экземпляр Storage с указанной корневой директорией.
func (s *Storage) RootDir() string {
	return s.rootDir
}

// SaveFile сохраняет данные в файл по указанному локальному пути.
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
