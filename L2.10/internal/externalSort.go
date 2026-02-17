package internal

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sort"
	"sync"
)

func Sort() {

}

type chunkResult struct {
	path string
	err  error
}

const bufSize = 1 << 20

func makeSortedChunks(r io.Reader, less LessFunc, workers int) ([]string, error) {

	chIn := make(chan []string, workers) // чанки строк
	resCh := make(chan chunkResult, workers)
	// работыши принимают значения из chIn, сортируют, отдают на запись в temp, и затем отправляют названия файлов в resCh
	wg := startChunkWorkers(chIn, resCh, less, workers)
	// закрыть resCh, когда воркеры закончат
	go func() {
		wg.Wait()
		close(resCh)
	}()

	lengthFile, err := lineCounter(r)
	if err != nil {
		return nil, err
	}

	// читает строки из файла, который нужно вывести, в chIn для работышей
	prodErr := produceChunks(r, chIn, lengthFile)
	// закрываем вход для работышей и ждём
	close(chIn)

	// собираем результаты (и удаляем мусор при ошибке)
	files, collectErr := collectChunkResults(resCh)

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

// работыши принимают массивы строк, сортируют, отдают на запись в temp, и затем отправляют названия файлов в resCh
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
	}
	if err := bw.Flush(); err != nil {
		return "", err
	}
	return file.Name(), nil
}

// читает строки из файла, который нужно вывести, в канал для работышей
func produceChunks(r io.Reader, chIn chan<- []string, lengthFile int) error {
	br := bufio.NewReaderSize(r, bufSize)
	capCurLines := min(lengthFile/8, bufSize)
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

		curLines = append(curLines, line)
		curBytes += int64(len(line))

		if curBytes >= bufSize {
			flush()
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

// подсчет количества строчек в исходном файле
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32<<10)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
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
