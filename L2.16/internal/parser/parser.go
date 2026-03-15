package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"wget/internal/urlutil"

	"golang.org/x/net/html"
)

type LinkKind string

const (
	LinkPage       LinkKind = "page"
	LinkStylesheet LinkKind = "stylesheet"
	LinkScript     LinkKind = "script"
	LinkImage      LinkKind = "image"
)

type Link struct {
	URL  string
	Kind LinkKind
}

func IsHTML(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "text/html")
}

func ExtractLinks(htmlData []byte, baseURL string) ([]Link, error) {
	root, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}

	var links []Link
	walkHTML(root, parsedBaseURL, &links)

	return uniqueLinks(links), nil
}

func walkHTML(node *html.Node, baseURL *url.URL, links *[]Link) {
	if node.Type == html.ElementNode {
		switch node.Data {
		case "a":
			if href, ok := getAttr(node, "href"); ok {
				if absURL, ok := resolveURL(baseURL, href); ok {
					*links = append(*links, Link{
						URL:  absURL,
						Kind: LinkPage,
					})
				}
			}
		case "link":
			if href, ok := getAttr(node, "href"); ok {
				if absURL, ok := resolveURL(baseURL, href); ok {
					*links = append(*links, Link{
						URL:  absURL,
						Kind: LinkStylesheet,
					})
				}
			}
		case "script":
			if src, ok := getAttr(node, "src"); ok {
				if absURL, ok := resolveURL(baseURL, src); ok {
					*links = append(*links, Link{
						URL:  absURL,
						Kind: LinkScript,
					})
				}
			}
		case "img":
			if src, ok := getAttr(node, "src"); ok {
				if absURL, ok := resolveURL(baseURL, src); ok {
					*links = append(*links, Link{
						URL:  absURL,
						Kind: LinkImage,
					})
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		walkHTML(child, baseURL, links)
	}
}

func getAttr(node *html.Node, key string) (string, bool) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}

func resolveURL(baseURL *url.URL, rawLink string) (string, bool) {
	rawLink = strings.TrimSpace(rawLink)
	if rawLink == "" {
		return "", false
	}

	lowerLink := strings.ToLower(rawLink)

	if strings.HasPrefix(rawLink, "#") {
		return "", false
	}
	if strings.HasPrefix(lowerLink, "javascript:") {
		return "", false
	}
	if strings.HasPrefix(lowerLink, "mailto:") {
		return "", false
	}
	if strings.HasPrefix(lowerLink, "tel:") {
		return "", false
	}

	parsedLink, err := url.Parse(rawLink)
	if err != nil {
		return "", false
	}

	resolved := baseURL.ResolveReference(parsedLink)

	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return "", false
	}

	return resolved.String(), true
}

func uniqueLinks(items []Link) []Link {
	seen := make(map[string]struct{})
	result := make([]Link, 0, len(items))

	for _, item := range items {
		key := string(item.Kind) + "|" + item.URL

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, item)
	}

	return result
}

func RewriteHTMLLinks(htmlData []byte, pageURL string, pageLocalPath string, rootDir string) ([]byte, error) {
	rootNode, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return nil, fmt.Errorf("parse html for rewrite: %w", err)
	}

	parsedPageURL, err := url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("parse page url: %w", err)
	}

	rewriteHTMLTree(rootNode, parsedPageURL, pageLocalPath, rootDir)

	var buf bytes.Buffer
	if err := html.Render(&buf, rootNode); err != nil {
		return nil, fmt.Errorf("render rewritten html: %w", err)
	}

	return buf.Bytes(), nil
}
func rewriteHTMLTree(node *html.Node, pageURL *url.URL, pageLocalPath string, rootDir string) {
	if node.Type == html.ElementNode {
		switch node.Data {
		case "a", "link":
			rewriteAttr(node, "href", pageURL, pageLocalPath, rootDir)
		case "script", "img":
			rewriteAttr(node, "src", pageURL, pageLocalPath, rootDir)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		rewriteHTMLTree(child, pageURL, pageLocalPath, rootDir)
	}
}

func rewriteAttr(node *html.Node, attrKey string, pageURL *url.URL, pageLocalPath string, rootDir string) {
	for i, attr := range node.Attr {
		if attr.Key != attrKey {
			continue
		}

		resolvedURL, ok := resolveURL(pageURL, attr.Val)
		if !ok {
			continue
		}

		localRelPath, ok := localRelativePath(resolvedURL, pageLocalPath, rootDir)
		if !ok {
			continue
		}

		node.Attr[i].Val = localRelPath
	}
}

func localRelativePath(targetURL string, currentPageLocalPath string, rootDir string) (string, bool) {
	targetLocalPath, err := urlutil.LocalPathForRewrite(rootDir, targetURL)
	if err != nil {
		return "", false
	}

	currentDir := filepath.Dir(currentPageLocalPath)

	relPath, err := filepath.Rel(currentDir, targetLocalPath)
	if err != nil {
		return "", false
	}

	return filepath.ToSlash(relPath), true
}
