package urlutil

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

func LocalPath(rootDir string, rawURL string, contentType string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf("url has empty host: %s", rawURL)
	}

	cleanPath := path.Clean(parsedURL.Path)

	if cleanPath == "." || cleanPath == "/" {
		return filepath.Join(rootDir, parsedURL.Host, "index.html"), nil
	}

	relativePath := filepath.FromSlash(strings.TrimPrefix(cleanPath, "/"))

	if strings.HasSuffix(parsedURL.Path, "/") {
		return filepath.Join(rootDir, parsedURL.Host, relativePath, "index.html"), nil
	}

	if hasExtension(path.Base(cleanPath)) {
		return filepath.Join(rootDir, parsedURL.Host, relativePath), nil
	}

	if isHTMLContentType(contentType) {
		return filepath.Join(rootDir, parsedURL.Host, relativePath, "index.html"), nil
	}
	ext := extensionByContentType(contentType)
	if ext != "" {
		return filepath.Join(rootDir, parsedURL.Host, relativePath+ext), nil
	}
	return filepath.Join(rootDir, parsedURL.Host, relativePath), nil
}

func hasExtension(name string) bool {
	ext := path.Ext(name)
	return ext != ""
}

func isHTMLContentType(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "text/html")
}
func extensionByContentType(contentType string) string {
	ct := strings.ToLower(contentType)

	switch {
	case strings.Contains(ct, "image/jpeg"):
		return ".jpg"
	case strings.Contains(ct, "image/png"):
		return ".png"
	case strings.Contains(ct, "image/gif"):
		return ".gif"
	case strings.Contains(ct, "text/css"):
		return ".css"
	case strings.Contains(ct, "javascript"):
		return ".js"
	case strings.Contains(ct, "image/webp"):
		return ".webp"
	case strings.Contains(ct, "image/svg+xml"):
		return ".svg"
	case strings.Contains(ct, "application/json"):
		return ".json"
	case strings.Contains(ct, "text/plain"):
		return ".txt"
	default:
		return ""
	}
}

func LocalPathForRewrite(rootDir string, rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}

	if parsedURL.Host == "" {
		return "", fmt.Errorf("url has empty host: %s", rawURL)
	}

	cleanPath := path.Clean(parsedURL.Path)

	if cleanPath == "." || cleanPath == "/" {
		return filepath.Join(rootDir, parsedURL.Host, "index.html"), nil
	}

	relativePath := filepath.FromSlash(strings.TrimPrefix(cleanPath, "/"))

	if strings.HasSuffix(parsedURL.Path, "/") {
		return filepath.Join(rootDir, parsedURL.Host, relativePath, "index.html"), nil
	}

	if hasExtension(path.Base(cleanPath)) {
		return filepath.Join(rootDir, parsedURL.Host, relativePath), nil
	}

	return filepath.Join(rootDir, parsedURL.Host, relativePath, "index.html"), nil
}
