package crawler

import (
	"fmt"
	"net/url"
	"strings"

	"wget/internal/fetcher"
	"wget/internal/parser"
	"wget/internal/storage"
	"wget/internal/urlutil"
)

type Crawler struct {
	fetcher   *fetcher.Fetcher
	storage   *storage.Storage
	visited   map[string]struct{}
	host      string
	entryPath string
}

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

func (c *Crawler) EntryPath() string {
	return c.entryPath
}

func (c *Crawler) CrawlPage(rawURL string, depth int) error {
	if depth < 0 {
		return nil
	}

	if !c.isSameHost(rawURL) {
		return nil
	}

	if c.isVisited(rawURL) {
		return nil
	}

	c.markVisited(rawURL)

	fmt.Printf("download page: %s (depth=%d)\n", rawURL, depth)

	resp, err := c.fetcher.Fetch(rawURL)
	if err != nil {
		return fmt.Errorf("fetch page %s: %w", rawURL, err)
	}

	localPath, err := urlutil.LocalPath(c.storage.RootDir(), resp.FinalURL, resp.ContentType)
	if err != nil {
		return fmt.Errorf("build local path for %s: %w", resp.FinalURL, err)
	}
	if c.entryPath == "" {
		c.entryPath = localPath
	}
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
			if depth > 0 {
				if !c.isSameHost(link.URL) && link.Kind == parser.LinkPage {
					fmt.Printf("skip external host: [%s] %s\n", link.Kind, link.URL)
					continue
				}

				if err := c.CrawlPage(link.URL, depth-1); err != nil {
					fmt.Printf("warn: %v\n", err)
				}
			}

		case parser.LinkStylesheet, parser.LinkScript, parser.LinkImage:
			if !c.isAllowedResourceHost(link.URL) {
				fmt.Printf("skip external resource: [%s] %s\n", link.Kind, link.URL)
				continue
			}
			if err := c.DownloadResource(link.URL); err != nil {
				fmt.Printf("warn: %v\n", err)
			}
		}
	}
	return nil
}
func (c *Crawler) DownloadResource(rawURL string) error {

	if c.isVisited(rawURL) {
		return nil
	}

	c.markVisited(rawURL)

	fmt.Printf("download resource: %s\n", rawURL)

	resp, err := c.fetcher.Fetch(rawURL)
	if err != nil {
		return fmt.Errorf("fetch resource %s: %w", rawURL, err)
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

func (c *Crawler) isVisited(rawURL string) bool {
	_, exists := c.visited[rawURL]
	return exists
}

func (c *Crawler) markVisited(rawURL string) {
	c.visited[rawURL] = struct{}{}
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
