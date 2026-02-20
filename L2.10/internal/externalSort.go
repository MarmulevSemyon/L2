package internal

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

func Sort() {

}

type chunkResult struct {
	path string
	err  error
}

const bufSize = 1 << 20
const chunkLimitBytes = 64 << 20 // лимит накопления строк в памяти перед сортировкой

// MakeSortedChunks cоздаёт набор файлов с отсортированными строками.
// r - ридер соритруемого файла;
// less - функция сравнения двух строк;
// workers - число работышей.
// возвращает массив названий файлов и ошибку.
func MakeSortedChunks(r io.Reader, less LessFunc, workers int) ([]string, error) {

	chIn := make(chan []string, workers) // чанки строк
	resCh := make(chan chunkResult, workers)
	// работыши принимают значения из chIn, сортируют, отдают на запись в temp, и затем отправляют названия файлов в resCh
	wg := startChunkWorkers(chIn, resCh, less, workers)
	// закрыываем resCh, когда воркеры закончат
	go func() {
		wg.Wait()
		close(resCh)
	}()

	// читаем строки из файла, который нужно вывести, в chIn для работышей
	prodErr := produceChunks(r, chIn)
	// закрываем вход для работышей и ждём
	close(chIn)

	// собираем результаты
	files, collectErr := collectChunkResults(resCh)
	// и удаляем мусор при ошибке
	if prodErr != nil {
		cleanupFiles(files)
		return nil, prodErr
	}
	if collectErr != nil {
		cleanupFiles(files)
		return nil, collectErr
	}

	return files, nil
}

// работыши принимают массивы строк, сортируют, отдают на запись в tempFile, и затем отправляют названия файлов в resCh
func startChunkWorkers(chIn <-chan []string, resCh chan<- chunkResult, less LessFunc, workers int) *sync.WaitGroup {
	var wg = sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for lines := range chIn {
				sort.Slice(lines, func(i, j int) bool {
					return less(lines[i], lines[j])
				})
				path, err := writeChunkToTempFile(lines)
				resCh <- chunkResult{path: path, err: err}
			}
		}()
	}
	return &wg
}

// принимает от работыша отсортированный массив, пишут его в файл и возвращают его название
func writeChunkToTempFile(lines []string) (string, error) {
	file, err := os.CreateTemp("", "sorted-chunck-*.txt")
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	bw := bufio.NewWriterSize(file, bufSize)
	for _, s := range lines {
		if _, err := bw.WriteString(s); err != nil {
			return "", err
		}
		if err := bw.WriteByte('\n'); err != nil {
			return "", err
		}
	}
	if err := bw.Flush(); err != nil {
		return "", err
	}
	return file.Name(), nil
}

// читает строки из исходного файла в канал для работышей
func produceChunks(r io.Reader, chIn chan<- []string) error {
	br := bufio.NewReaderSize(r, bufSize)
	capCurLines := bufSize / 8
	curLines := make([]string, 0, capCurLines)
	var curBytes int64

	flush := func() {
		if len(curLines) == 0 {
			return
		}
		chIn <- curLines
		curLines = nil
		curBytes = 0
	}
	for {
		line, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		line = strings.TrimRight(line, "\r\n")
		curLines = append(curLines, line)
		curBytes += int64(len(line))

		if curBytes >= chunkLimitBytes {
			flush()
		}
		if err == io.EOF {
			break
		}
	}
	flush()
	return nil
}

// собирает названия файлов от работышей
func collectChunkResults(resCh <-chan chunkResult) ([]string, error) {
	files := make([]string, 0, 10)
	for res := range resCh {
		if res.err != nil {
			return files, res.err
		}
		files = append(files, res.path)
	}
	return files, nil
}

func cleanupFiles(paths []string) {
	for _, p := range paths {
		_ = os.Remove(p)
	}
}
