package crawler

import (
	"fmt"

	"gosearch/models"
	"gosearch/pkg/parser"
	"gosearch/internal/repositories"
)

type Crawler interface {
	Scan(url string, depth int) (string, error)
	BatchScan(urls []string, depth int, workers int) (<-chan models.Document, <-chan error)
}

type crawlerImpl struct {
	parser parser.Parser
	repo repositories.DocumentRepository
}

func New() *crawlerImpl {
	return &crawlerImpl{
		parser: parser.New(),
		repo: repositories.NewDocumentRepository(),
	}
}

func (c *crawlerImpl) Scan(url string, depth int) error {
	if depth <= 0 {
		return nil
	}

	resp, err := c.parser.FetchHTML(url)
	if err != nil {
		return fmt.Errorf("error reading URL %s > %w", url, err)
	}

	links := c.parser.ExtractLinks(resp)
	groupedLinks := c.parser.GroupLinks(links)

	for _, link := range groupedLinks {
		if err := c.repo.FindOrCreate(resp, map[string]any{"lg": &link}, depth); err != nil {
			return err
		}

		if err := c.Scan(link.Url, depth-1); err != nil {
			return err
		}
	}

	return nil
}
