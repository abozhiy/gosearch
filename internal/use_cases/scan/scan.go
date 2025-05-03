package scan

import (
	// "os"

	"gosearch/pkg/crawler"
	// "gosearch/pkg/searcher"
)



func Perform(url string) {
	// result, err := searcher.FindURLData(url)
	// if err != nil {
	// 	panic(err)
	// }

	// if result {
	// 	return result
	// }

	c := crawler.New()
	if err := c.Scan(url, 2); err != nil {
		panic(err)
	}
}
