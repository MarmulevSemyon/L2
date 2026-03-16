package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		args        []string
		expected    Config
		wantErr     bool
		errContains string
	}{
		{
			name: "успешный разбор с таймаутом по умолчанию",
			args: []string{"localhost", "8080"},
			expected: Config{
				Timeout: 10 * time.Second,
				Host:    "localhost",
				Port:    "8080",
			},
		},
		{
			name: "успешный разбор с длинным флагом таймаута",
			args: []string{"--timeout=5s", "example.com", "23"},
			expected: Config{
				Timeout: 5 * time.Second,
				Host:    "example.com",
				Port:    "23",
			},
		},
		{
			name: "успешный разбор с коротким флагом таймаута",
			args: []string{"-T", "1m", "127.0.0.1", "9000"},
			expected: Config{
				Timeout: 1 * time.Minute,
				Host:    "127.0.0.1",
				Port:    "9000",
			},
		},
		{
			name:        "ошибка если аргументы не переданы",
			args:        []string{},
			wantErr:     true,
			errContains: "usage",
		},
		{
			name:        "ошибка если передан только хост",
			args:        []string{"localhost"},
			wantErr:     true,
			errContains: "usage",
		},
		{
			name:        "ошибка если передано слишком много позиционных аргументов",
			args:        []string{"localhost", "8080", "extra"},
			wantErr:     true,
			errContains: "usage",
		},
		{
			name:        "ошибка если таймаут некорректный",
			args:        []string{"--timeout=abc", "localhost", "8080"},
			wantErr:     true,
			errContains: "parse flags",
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(testCase.args)

			if testCase.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, testCase.expected, got)
		})
	}
}
