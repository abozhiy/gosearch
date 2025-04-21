package searcher

import (

)

type SearchInterface interface {
	FindURLData(url string) (map[string]string, error)
}