package crawler

import (
	"fmt"

	"gosearch/models"
	"gosearch/pkg/parser"
	"gosearch/internal/repositories"
)

type CrawlerInterface interface {
	Scan(url string, depth int) (string, error)
	BatchScan(urls []string, depth int, workers int) (<-chan models.Document, <-chan error)
}

func Scan(url string, depth int) error {
	if depth <= 0 {
		return nil
	}

	resp, err := parser.ReadUrl(url)
	if err != nil {
		return fmt.Errorf("error reading URL %s > %w", url, err)
	}

	links := parser.ExtractLinks(resp)
	groupedLinks := parser.GroupLinks(links)

	for _, link := range groupedLinks {
		if err := document_repository.FindOrCreate(resp, &link, depth); err != nil {
			return err
		}

		if err := Scan(link.Url, depth-1); err != nil {
			return err
		}
	}

	return nil
}
