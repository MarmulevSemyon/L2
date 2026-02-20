package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintIsSortedLexic(t *testing.T) {
	// test -c
	flags := LineArgs{
		C: true,
	}
	less, _ := BuildLess(flags)

	file := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test.txt"
	response, err := PrintIsSorted(file, less)
	assert.Nil(t, err)
	assert.Equal(t, "Файл отсортирован!", response)

	//test1 -сr
	flags1 := LineArgs{
		C: true,
		R: true,
	}
	less1, _ := BuildLess(flags1)

	file1 := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test1.txt"
	response1, err1 := PrintIsSorted(file1, less1)
	assert.Nil(t, err1)
	assert.Equal(t, "Файл отсортирован!", response1)

	// test2 -c -b
	flags2 := LineArgs{
		C: true,
		B: true,
	}
	less2, _ := BuildLess(flags2)

	file2 := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test2.txt"
	response2, err2 := PrintIsSorted(file2, less2)
	assert.Nil(t, err2)
	assert.Equal(t, "Файл отсортирован!", response2)

	//test -c -b -r
	flags3 := LineArgs{
		C: true,
		R: true,
		B: true,
	}
	less3, _ := BuildLess(flags3)

	file3 := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test3.txt"
	response3, err3 := PrintIsSorted(file3, less3)
	assert.Nil(t, err3)
	assert.Equal(t, "Файл отсортирован!", response3)
}

func TestPrintIsSortedMonth(t *testing.T) {
	// test -cr
	flags := LineArgs{
		C: true,
		R: true,
	}
	less, _ := BuildLess(flags)

	file := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test.txt"
	response, err := PrintIsSorted(file, less)
	assert.Nil(t, err)
	assert.Equal(t, "Файл отсортирован!", response)

	//test1 -с
	flags1 := LineArgs{
		C: true,
	}
	less1, _ := BuildLess(flags1)

	file1 := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test1.txt"
	response1, err1 := PrintIsSorted(file1, less1)
	assert.Nil(t, err1)
	assert.Equal(t, "Файл отсортирован!", response1)

	// test2 -c -b
	flags2 := LineArgs{
		C: true,
		B: true,
	}
	less2, _ := BuildLess(flags2)

	file2 := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test2.txt"
	response2, err2 := PrintIsSorted(file2, less2)
	assert.Nil(t, err2)
	assert.Equal(t, "Файл отсортирован!", response2)

	//test -c -b -r
	flags3 := LineArgs{
		C: true,
		R: true,
		B: true,
	}
	less3, _ := BuildLess(flags3)

	file3 := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test3.txt"
	response3, err3 := PrintIsSorted(file3, less3)
	assert.Nil(t, err3)
	assert.Equal(t, "Файл отсортирован!", response3)
}
