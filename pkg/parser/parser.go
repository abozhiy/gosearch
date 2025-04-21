package parser

import (
	"fmt"
	"io"
	"strings"

	"net/http"
	"net/url"

	"golang.org/x/net/html"

	"gosearch/internal/types"
)

type DocumentParser interface {
	ReadUrl(url string) ([]uint8, error)
	ExtractLinks(resp []uint8) []types.Link
	GroupLinks(links []types.Link) []types.LinkGrouped
}

func ReadUrl(url string) ([]uint8, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("while geting request (%v)  > %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading response data > %w", err)
	}

	return body, nil
}

func ExtractLinks(resp []uint8) []types.Link {
	var links []types.Link

	doc, err := html.Parse(strings.NewReader(string(resp)))
	if err != nil {
		return links
	}

	var crawl func(*html.Node)
	crawl = func(doc *html.Node) {
		if doc.Type == html.ElementNode && doc.Data == "a" {
			if href := extractHref(doc); href != "" {
				if title := extractText(doc); !strings.Contains(title, "<img") {
					link := types.Link{Url: href, Title: title}
					links = append(links, link)
				}
			}
		}

		for c := doc.FirstChild; c != nil; c = c.NextSibling {
			crawl(c)
		}
	}
	crawl(doc)

	return links
}

func GroupLinks(links []types.Link) []types.LinkGrouped {
	urlTitlesSet := make(map[string]map[string]bool)

	for _, link := range links {
		if urlTitlesSet[link.Url] == nil {
			urlTitlesSet[link.Url] = make(map[string]bool)
		}
		urlTitlesSet[link.Url][link.Title] = true
	}

	var result []types.LinkGrouped
	for url, titlesSet := range urlTitlesSet {
		var titles []string
		for title := range titlesSet {
			titles = append(titles, title)
		}

		if len(titles) == 0 {
			continue
		}

		result = append(result, types.LinkGrouped{Url: url, Titles: titles})
	}

	return result
}

func extractHref(n *html.Node) string {
	href := getAttrValue(n, "href")
	if href == "" {
		return ""
	}

	if hasSrcElement(n) {
		return ""
	}

	parsed, err := url.Parse(strings.TrimSpace(href))
	if err != nil {
		return ""
	}

	if parsed.Scheme != "https" || parsed.Host == "" {
		return ""
	}

	return stripUrl(href)
}

func hasSrcElement(n *html.Node) bool {
	if n.FirstChild == nil {
		return false
	}

	for _, attr := range n.FirstChild.Attr {
		if attr.Key == "src" {
			return true
		}
	}

	return false
}

func getAttrValue(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}

	var sb strings.Builder

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(extractText(c))
	}

	return sb.String()
}

func stripUrl(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	parsed.RawQuery = ""
	return parsed.String()
}
