package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"wget/internal/fetcher"
	"wget/internal/parser"
	"wget/internal/storage"
	"wget/internal/urlutil"
)

// Task описывает одну задачу на обработку URL.
type Task struct {
	URL   string
	Depth int
	Kind  parser.LinkKind
}

// Crawler выполняет обход сайта, скачивание страниц и ресурсов.
type Crawler struct {
	fetcher *fetcher.Fetcher
	storage *storage.Storage

	host string

	visited map[string]struct{}
	mu      sync.Mutex

	entryPath string
}

// Crawler выполняет обход сайта, скачивание страниц и ресурсов.
func New(fetcher *fetcher.Fetcher, storage *storage.Storage, startURL string) (*Crawler, error) {
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("parse start url: %w", err)
	}

	if parsedURL.Host == "" {
		return nil, fmt.Errorf("start url has empty host")
	}

	return &Crawler{
		fetcher: fetcher,
		storage: storage,
		visited: make(map[string]struct{}),
		host:    parsedURL.Host,
	}, nil
}

// New создаёт новый экземпляр Crawler для стартового URL.
func (c *Crawler) EntryPath() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.entryPath
}

func (c *Crawler) setEntryPath(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.entryPath == "" {
		c.entryPath = path
	}
}

func (c *Crawler) markVisitedIfNew(rawURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.visited[rawURL]; exists {
		return false
	}

	c.visited[rawURL] = struct{}{}
	return true
}

// Run запускает приложение и координирует процесс скачивания сайта.
func (c *Crawler) Run(startURL string, depth int, concurrency int) error {

	tasks := make(chan Task, concurrency*4)

	var workersWG sync.WaitGroup
	var tasksWG sync.WaitGroup

	worker := func() {
		defer workersWG.Done()

		for task := range tasks {
			c.processTask(task, tasks, &tasksWG)
			tasksWG.Done()
		}
	}

	for i := 0; i < concurrency; i++ {
		workersWG.Add(1)
		go worker()
	}

	tasksWG.Add(1)
	tasks <- Task{
		URL:   startURL,
		Depth: depth,
		Kind:  parser.LinkPage,
	}

	tasksWG.Wait()
	close(tasks)
	workersWG.Wait()

	return nil
}

func (c *Crawler) processTask(task Task, tasks chan<- Task, tasksWG *sync.WaitGroup) {
	switch task.Kind {
	case parser.LinkPage:
		if err := c.processPage(task, tasks, tasksWG); err != nil {
			fmt.Printf("warn: %v\n", err)
		}
	case parser.LinkStylesheet, parser.LinkScript, parser.LinkImage:
		if err := c.processResource(task); err != nil {
			fmt.Printf("warn: %v\n", err)
		}
	default:
		fmt.Printf("warn: unknown task kind for %s\n", task.URL)
	}
}

func (c *Crawler) processPage(task Task, tasks chan<- Task, tasksWG *sync.WaitGroup) error {
	if task.Depth < 0 {
		return nil
	}

	if !c.isSameHost(task.URL) {
		return nil
	}

	if !c.markVisitedIfNew(task.URL) {
		return nil
	}

	fmt.Printf("download page: %s (depth=%d)\n", task.URL, task.Depth)

	resp, err := c.fetcher.Fetch(task.URL)
	if err != nil {
		return fmt.Errorf("fetch page %s: %w", task.URL, err)
	}

	localPath, err := urlutil.LocalPath(c.storage.RootDir(), resp.FinalURL, resp.ContentType)
	if err != nil {
		return fmt.Errorf("build local path for %s: %w", resp.FinalURL, err)
	}

	c.setEntryPath(localPath)

	dataToSave := resp.Body
	if parser.IsHTML(resp.ContentType) {
		rewrittenHTML, err := parser.RewriteHTMLLinks(resp.Body, resp.FinalURL, localPath, c.storage.RootDir())
		if err != nil {
			fmt.Printf("warn: rewrite html %s: %v\n", resp.FinalURL, err)
		} else {
			dataToSave = rewrittenHTML
		}
	}

	if err := c.storage.SaveFile(localPath, dataToSave); err != nil {
		return fmt.Errorf("save page %s: %w", localPath, err)
	}

	fmt.Printf("saved: %s\n", localPath)

	if !parser.IsHTML(resp.ContentType) {
		return nil
	}

	links, err := parser.ExtractLinks(resp.Body, resp.FinalURL)
	if err != nil {
		return fmt.Errorf("extract links from %s: %w", resp.FinalURL, err)
	}

	for _, link := range links {
		fmt.Printf("found [%s]: %s\n", link.Kind, link.URL)

		switch link.Kind {
		case parser.LinkPage:
			if task.Depth > 0 && c.isSameHost(link.URL) {
				// кладём в отдельной горутине
				enqueueTask(tasks, tasksWG, Task{
					URL:   link.URL,
					Depth: task.Depth - 1,
					Kind:  parser.LinkPage,
				})
			}

		case parser.LinkStylesheet, parser.LinkScript, parser.LinkImage:
			if c.isAllowedResourceHost(link.URL) {
				// кладём в отдельной горутине
				enqueueTask(tasks, tasksWG, Task{
					URL:   link.URL,
					Depth: task.Depth,
					Kind:  link.Kind,
				})
			} else {
				fmt.Printf("skip external resource: [%s] %s\n", link.Kind, link.URL)
			}
		}
	}

	return nil
}

func (c *Crawler) processResource(task Task) error {
	if !c.isAllowedResourceHost(task.URL) {
		return nil
	}

	if !c.markVisitedIfNew(task.URL) {
		return nil
	}

	fmt.Printf("download resource: %s\n", task.URL)

	resp, err := c.fetcher.Fetch(task.URL)
	if err != nil {
		return fmt.Errorf("fetch resource %s: %w", task.URL, err)
	}

	localPath, err := urlutil.LocalPath(c.storage.RootDir(), resp.FinalURL, resp.ContentType)
	if err != nil {
		return fmt.Errorf("build local path for resource %s: %w", resp.FinalURL, err)
	}

	if err := c.storage.SaveFile(localPath, resp.Body); err != nil {
		return fmt.Errorf("save resource %s: %w", localPath, err)
	}

	fmt.Printf("saved: %s\n", localPath)

	return nil
}

func (c *Crawler) isSameHost(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return parsedURL.Host == c.host
}

func (c *Crawler) isAllowedResourceHost(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if parsedURL.Host == c.host {
		return true
	}

	return strings.HasSuffix(parsedURL.Host, "."+c.host)
}

// если класть не в отдельной горутине, то работыш заблочится при заполнении буфераи будет deadlock
func enqueueTask(tasks chan<- Task, tasksWG *sync.WaitGroup, task Task) {
	tasksWG.Add(1)

	go func() {
		tasks <- task
	}()
}
