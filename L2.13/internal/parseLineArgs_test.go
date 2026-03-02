package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateArgsNoFile(t *testing.T) {
	t.Parallel()

	err := validateArgs(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "не указан файл")
}

func TestParseLineNoFileArg(t *testing.T) {
	t.Parallel()

	_, _, err := ParseLine([]string{"-f", "1,3-5"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "не указан файл")
}

func TestParseLineFDSFile(t *testing.T) {
	t.Parallel()

	flags, rest, err := ParseLine([]string{"-f", "1,3-5", "-d", ",", "-s", "input.txt"})
	require.NoError(t, err)

	assert.Equal(t, "1,3-5", flags.Fields)
	assert.Equal(t, ",", flags.Delimiter)
	assert.True(t, flags.Separated)
	assert.Equal(t, []string{"input.txt"}, rest)
}

func TestParseLine_DefaultDTab(t *testing.T) {
	t.Parallel()

	flags, rest, err := ParseLine([]string{"-f", "1", "file.tsv"})
	require.NoError(t, err)

	assert.Equal(t, "\t", flags.Delimiter)
	assert.Equal(t, []string{"file.tsv"}, rest)
}

func TestParseLine_ProfileFlags(t *testing.T) {
	t.Parallel()

	flags, rest, err := ParseLine([]string{
		"--cpuprofile", "cpu.out",
		"--memprofile", "mem.out",
		"--trace", "trace.out",
		"file.txt",
	})
	require.NoError(t, err)

	assert.Equal(t, "cpu.out", flags.CPUProfile)
	assert.Equal(t, "mem.out", flags.MemProfile)
	assert.Equal(t, "trace.out", flags.Trace)
	assert.Equal(t, []string{"file.txt"}, rest)
}

func TestParseLine_UnknownFlagReturnsError(t *testing.T) {
	t.Parallel()

	_, _, err := ParseLine([]string{"--no-such-flag", "file.txt"})
	require.Error(t, err)
}
