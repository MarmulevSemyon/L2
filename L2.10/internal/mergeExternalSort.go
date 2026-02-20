package internal

import (
	"bufio"
	"container/heap"
	"io"
	"os"
)

type heapItem struct {
	line    string
	fileIdx int
}

type mergeHeap struct {
	items []heapItem
	less  LessFunc
}

func (h mergeHeap) Len() int { return len(h.items) }
func (h mergeHeap) Less(i, j int) bool {
	return h.less(h.items[i].line, h.items[j].line)
}
func (h mergeHeap) Swap(i, j int) { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *mergeHeap) Push(x any) {
	h.items = append(h.items, x.(heapItem))
}

func (h *mergeHeap) Pop() any {
	n := len(h.items)
	x := h.items[n-1]
	h.items = h.items[:n-1]
	return x
}

// MergeSortedFilesToWriterHeap сливает отсортированные файлы (paths) в writer.
// less — компаратор строк.
// unique — если true, убирает дубли (как -u) на лету.
func MergeSortedFilesToWriterHeap(paths []string, w io.Writer, less LessFunc, unique bool) error {
	type fileState struct {
		f  *os.File
		br *bufio.Reader
	}

	// открываем файлы и кладём в []fileState
	states := make([]fileState, 0, len(paths))
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			for _, st := range states {
				_ = st.f.Close()
			}
			return err
		}
		states = append(states, fileState{
			f:  f,
			br: bufio.NewReaderSize(f, bufSize),
		})
	}
	defer func() {
		for _, st := range states {
			_ = st.f.Close()
		}
	}()

	// создаем бинарную кучу
	h := &mergeHeap{less: less}
	heap.Init(h)

	for i := range states {
		line, ok, err := readLineKeepNL(states[i].br)
		if err != nil {
			return err
		}
		if ok {
			heap.Push(h, heapItem{line: line, fileIdx: i})
		}
	}
	// создаём писателя через буфер в stdout
	bw := bufio.NewWriterSize(w, bufSize)
	defer bw.Flush()

	var prev string
	hasPrev := false

	//  Основной слив
	for h.Len() > 0 {
		it := heap.Pop(h).(heapItem) // утверждение типа, type assertion
		//если нет -u, то пишем
		if !unique || !hasPrev || !equalByLess(prev, it.line, less) {
			if _, err := bw.WriteString(it.line); err != nil {
				return err
			}
			prev = it.line
			hasPrev = true
		}

		// взять следующую строку из того же файла и положить обратно в кучу
		next, ok, err := readLineKeepNL(states[it.fileIdx].br)
		if err != nil {
			return err
		}
		if ok {
			heap.Push(h, heapItem{line: next, fileIdx: it.fileIdx})
		}
	}

	return nil
}

// читает строку
func readLineKeepNL(br *bufio.Reader) (string, bool, error) {
	line, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", false, err
	}
	if len(line) == 0 && err == io.EOF {
		return "", false, nil
	}
	return line, true, nil
}

func equalByLess(a, b string, less LessFunc) bool {
	return !less(a, b) && !less(b, a)
}
