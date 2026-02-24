package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	lineArgsErr := [][]string{{"sort", "-c", "-r", "-n", "-u", "-M", "-b", "-h", "-k", "12"},
		{"sort", "-crnuMbhk12"},
		{"sort", "-crnuMbhk=12"},
		{"sort", "-crnuMbhk", "12"},
		{"sort", "-crnuMbh", "--key", "12"},
		{"sort", "--key", "12", "-crnuMbh"},
	}

	empty := LineArgs{}
	for i := range lineArgsErr {
		actual, _, err := ParseLine(lineArgsErr[i])
		assert.NotNil(t, err)
		assert.Equal(t, empty, actual)
	}

	lineArgsN := [][]string{{"sort", "-c", "-r", "-n", "-u", "-b", "-k", "12"},
		{"sort", "-crnubk12"},
		{"sort", "-crnubk=12"},
		{"sort", "-crnubk", "12"},
		{"sort", "-crnub", "--key", "12"},
		{"sort", "--key", "12", "-crnub"},
	}
	expectedN := LineArgs{
		K: 12,
		N: true,
		R: true,
		U: true,
		M: false,
		B: true,
		C: true,
		H: false,
	}
	for i := range lineArgsN {
		actual, _, err := ParseLine(lineArgsN[i])
		assert.Nil(t, err)
		assert.Equal(t, expectedN, actual)
	}

	lineArgsM := [][]string{{"sort", "-c", "-r", "-M", "-u", "-b", "-k", "12"},
		{"sort", "-crMubk12"},
		{"sort", "-crMubk=12"},
		{"sort", "-crMubk", "12"},
		{"sort", "-crMub", "--key", "12"},
		{"sort", "--key", "12", "-crMub"},
	}
	expectedM := LineArgs{
		K: 12,
		N: false,
		R: true,
		U: true,
		M: true,
		B: true,
		C: true,
		H: false,
	}
	for _, v := range lineArgsM {
		actual, _, err := ParseLine(v)
		assert.Nil(t, err)
		assert.Equal(t, expectedM, actual)
	}

	lineArgsH := [][]string{{"sort", "-c", "-r", "-h", "-u", "-b", "-k", "12"},
		{"sort", "-crhubk12"},
		{"sort", "-crhubk=12"},
		{"sort", "-crhubk", "12"},
		{"sort", "-crhub", "--key", "12"},
		{"sort", "--key", "12", "-crhub"},
	}
	expectedH := LineArgs{
		K: 12,
		N: false,
		R: true,
		U: true,
		M: false,
		B: true,
		C: true,
		H: true,
	}
	for _, v := range lineArgsH {
		actual, _, err := ParseLine(v)
		assert.Nil(t, err)
		assert.Equal(t, expectedH, actual)
	}

}
