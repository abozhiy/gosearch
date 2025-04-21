package parser_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gosearch/pkg/parser"
	"gosearch/internal/types"

	"github.com/stretchr/testify/require"
)

func TestReadUrl(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><body><a href="https://example.com">Example</a></body></html>`))
	}))
	defer server.Close()

	body, err := parser.ReadUrl(server.URL)

	require.NoError(t, err)
	require.Contains(t, string(body), "Example")
}

func TestExtractLinks(t *testing.T) {
	html := `
		<html>
		  <body>
		    <a href="https://example.com?label=gen173nr-1BEhVjb3ZpZF8xOV9ib2cfsdf3249raW5nX2Zh">Example1</a>
		    <a href="https://example.com?label=gen173nr-1BEhVjb3ZpZF8xOV91232xsadib29raW5nX2Zs">Example2</a>
		    <a href="https://another.com?label=gen173nr-1BEhVjb3ZpZF8xsweaew12312OV9ib29raW5nX">Another</a>
		    <a href="mailto:test@example.com">Email</a>
		  </body>
		</html>`

	links := parser.ExtractLinks([]byte(html))

	require.Len(t, links, 3)

	require.Equal(t, "https://example.com", links[0].Url)
	require.Equal(t, "Example1", links[0].Title)

	require.Equal(t, "https://example.com", links[1].Url)
	require.Equal(t, "Example2", links[1].Title)

	require.Equal(t, "https://another.com", links[2].Url)
	require.Equal(t, "Another", links[2].Title)
}

func TestExtractLinks_IgnoreImages(t *testing.T) {
	html := `<html><body>
	  <a href="https://img_url.com?label=gen173nr-1BEhVjb3ZpZF8xOV9ib2cfsdf3249raW5nX2Zh"><img src="x.jpg" />Click</a>
	  <a href="https://example.com?label=gen173nr-1BEhVjb3ZpZF8xOV91232xsadib29raW5nX2Zs">Real</a>
	</body></html>`

	links := parser.ExtractLinks([]byte(html))
	require.Len(t, links, 1)
	require.Equal(t, "https://example.com", links[0].Url)
	require.Equal(t, "Real", links[0].Title)
}

func TestGroupLinks(t *testing.T) {
	links := []types.Link{
		{Url: "https://example.com", Title: "Home"},
		{Url: "https://example.com", Title: "About"},
		{Url: "https://example.com", Title: "Home"},
		{Url: "https://other.com", Title: "Contact"},
	}

	expected := []types.LinkGrouped{
		{
			Url:    "https://example.com",
			Titles: []string{"Home", "About"},
		},
		{
			Url:    "https://other.com",
			Titles: []string{"Contact"},
		},
	}

	require.ElementsMatch(t, expected, parser.GroupLinks(links))
}
