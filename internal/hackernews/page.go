package hackernews

import (
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"net/http"
)

type Page struct {
	url       string
	bodyBytes []byte
}

func LoadPage(year, month, day int) (Page, error) {
	url := getPageURL(year, month, day)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Page{}, fmt.Errorf("failed to create request for page at %s: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Page{}, fmt.Errorf("error when making request for page %s: %w", url, err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Page{}, fmt.Errorf("error reading page %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("recieved non-200 status code %s: %d\n%s", url, resp.StatusCode, string(bodyBytes))
	}

	return Page{
		bodyBytes: bodyBytes,
	}, nil
}

func (p Page) ParseTitles() ([30]string, error) {
	root, err := htmlquery.Parse(bytes.NewReader(p.bodyBytes))
	if err != nil {
		return [30]string{}, fmt.Errorf("error parsing page %s: %ws", p.url, err)
	}

	titles, err := htmlquery.QueryAll(root, "//span[@class='titleline']/a/text()")
	if err != nil {
		return [30]string{}, fmt.Errorf("error querying paeg for titles %s: %w", p.url, err)
	}

	if len(titles) != 30 {
		return [30]string{}, fmt.Errorf("unexpected number of titles %s: %d", p.url, len(titles))
	}

	var result [30]string
	for i, title := range titles {
		if title.Data == "" {
			return [30]string{}, fmt.Errorf("title at index %d is empty %s", i, p.url)
		}

		result[i] = title.Data
	}

	return result, nil
}

func getPageURL(year, month, day int) string {
	return fmt.Sprintf("https://news.ycombinator.com/front?day=%d-%d-%d", year, month, day)
}
