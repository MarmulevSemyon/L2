package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintIsSortedLexic(t *testing.T) {
	// test -c
	flags := LineArgs{
		C: true,
	}
	less, _ := BuildLess(flags)
	fileName := filepath.Join("..", "test", "test.txt")
	file, _ := os.Open(fileName)
	defer file.Close()
	response, err := PrintIsSorted(file, less)
	assert.Nil(t, err)
	assert.Equal(t, "Файл отсортирован!", response)

	//test1 -сr
	flags1 := LineArgs{
		C: true,
		R: true,
	}
	less1, _ := BuildLess(flags1)
	fileName1 := filepath.Join("..", "test", "test1.txt")

	file1, _ := os.Open(fileName1)
	defer file1.Close()
	response1, err1 := PrintIsSorted(file1, less1)
	assert.Nil(t, err1)
	assert.Equal(t, "Файл отсортирован!", response1)

	// test2 -c -b
	flags2 := LineArgs{
		C: true,
		B: true,
	}
	less2, _ := BuildLess(flags2)

	fileName2 := filepath.Join("..", "test", "test2.txt")
	file2, _ := os.Open(fileName2)
	defer file2.Close()
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

	fileName3 := filepath.Join("..", "test", "test3.txt")
	file3, _ := os.Open(fileName3)
	defer file3.Close()
	response3, err3 := PrintIsSorted(file3, less3)
	assert.Nil(t, err3)
	assert.Equal(t, "Файл отсортирован!", response3)
}

func TestPrintIsSortedMonth(t *testing.T) {
	//test1 -сM
	flags := LineArgs{
		C: true,
		M: true,
	}
	less, _ := BuildLess(flags)

	fileName := filepath.Join("..", "test", "test_month.txt")
	file, _ := os.Open(fileName)
	defer file.Close()
	response, err := PrintIsSorted(file, less)
	assert.Nil(t, err)
	assert.Equal(t, "Файл отсортирован!", response)

	// test -cMr
	flags1 := LineArgs{
		C: true,
		M: true,
		R: true,
	}
	less1, _ := BuildLess(flags1)

	fileName1 := filepath.Join("..", "test", "test_month1.txt")
	file1, _ := os.Open(fileName1)
	defer file1.Close()
	response1, err1 := PrintIsSorted(file1, less1)
	assert.Nil(t, err1)
	assert.Equal(t, "Файл отсортирован!", response1)

	// // test2 -c -b
	// flags2 := LineArgs{
	// 	C: true,
	// 	B: true,
	// }
	// less2, _ := BuildLess(flags2)

	// fileName2 := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test2.txt"
	// file2, _ := os.Open(fileName2)
	// defer file2.Close()
	// response2, err2 := PrintIsSorted(file2, less2)
	// assert.Nil(t, err2)
	// assert.Equal(t, "Файл отсортирован!", response2)

	// //test -c -b -r
	// flags3 := LineArgs{
	// 	C: true,
	// 	R: true,
	// 	B: true,
	// }
	// less3, _ := BuildLess(flags3)

	// fileName3 := "C:\\Users\\79164\\Desktop\\proga\\GO\\WB\\L2\\L2.10\\test\\test3.txt"
	// file3, _ := os.Open(fileName3)
	// defer file3.Close()
	// response3, err3 := PrintIsSorted(file3, less3)
	// assert.Nil(t, err3)
	// assert.Equal(t, "Файл отсортирован!", response3)
}
