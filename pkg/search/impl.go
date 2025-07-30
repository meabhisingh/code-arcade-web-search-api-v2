package search

import (
	"bytes"
	"code-arcade-web-search-api/config"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	readability "github.com/go-shiori/go-readability"
)

func (s *service) FetchSearchResults(query string) ([]SearchResult, error) {

	searxngUrl := config.GetEnv("SEARCH_ENGINE_URL", "http://localhost:8080")

	urlWithQuery := searxngUrl + "?q=" + url.QueryEscape(query) + "&format=json"

	body := []byte(`{}`)

	resp, err := http.Post(urlWithQuery, "application/json", bytes.NewBuffer(body))

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()

	var data SearchEngineResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Results, nil
}

var contentSelectors = []string{
	"title",
	"meta[name='description']",
	"meta[property='og:description']",
	"article",
	"main",
	"section",
	"div",
	"[role='main']",
	".content",
	"#content",
	".post-content",
	".article-content",
	".entry-content",
	".main-content",
	"h1",
	"h2",
	"h3",
	"p",
	"blockquote",
	"ul",
	"ol",
	"li",
}

func (s *service) Scrape(results []SearchResult) ([]ScrapedPage, error) {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var pages []ScrapedPage

	client := &http.Client{Timeout: 10 * time.Second}

	for _, r := range results {
		wg.Add(1)
		go func(r SearchResult) {
			defer wg.Done()

			page, err := s.extractReadableContent(r.URL, client)
			if err != nil || page == nil || len(page.Content) < 200 {
				page, err = s.fallbackContent(r.URL, client)
			}

			if err != nil || page == nil {
				return
			}

			page.Content = cleanContent(page.Content)

			mu.Lock()
			pages = append(pages, *page)
			mu.Unlock()

		}(r)

	}
	wg.Wait()
	return pages, nil
}

func (s *service) extractReadableContent(rawURL string, client *http.Client) (*ScrapedPage, error) {
	req, _ := http.NewRequest(http.MethodGet, rawURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 ... Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return nil, err
	}

	scrapedPage := &ScrapedPage{
		Source: ScrapedPageSource{
			Url:     rawURL,
			Title:   article.Title,
			Website: article.SiteName,
			Icon:    article.Favicon,
		},
		Content: article.TextContent,
	}

	return scrapedPage, nil
}

func (s *service) fallbackContent(rawURL string, client *http.Client) (*ScrapedPage, error) {
	resp, err := client.Get(rawURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	for _, sel := range contentSelectors {
		text := strings.TrimSpace(doc.Find(sel).First().Text())
		if len(text) > 200 {
			return &ScrapedPage{
				Source:  ScrapedPageSource{Url: rawURL},
				Content: text,
			}, nil
		}
	}

	var paras []string
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if len(text) > 50 {
			paras = append(paras, text)
		}
	})
	return &ScrapedPage{
		Source:  ScrapedPageSource{Url: rawURL},
		Content: strings.Join(paras, "\n"),
	}, nil
}

func collapseRepeatedPunctuation(s string) string {
	var result strings.Builder
	prev := rune(0)
	count := 0

	for _, ch := range s {
		if strings.ContainsRune(".,!?", ch) && ch == prev {
			count++
			continue
		}
		if count > 0 {
			result.WriteRune(prev)
			count = 0
		}
		result.WriteRune(ch)
		prev = ch
	}
	if count > 0 {
		result.WriteRune(prev)
	}
	return result.String()
}

func cleanContent(text string) string {
	replacements := []struct {
		pattern string
		replace string
	}{
		{`\s+`, " "},
		{`https?://\S+`, ""},
		{`[\w\.-]+@[\w\.-]+\.\w+`, ""},
		{`[^\w\s.,!?-]`, " "},
		{`\s+([.,!?])`, "$1"},
		{`<!--[\s\S]*?-->`, ""},
	}

	for _, r := range replacements {
		text = regexp.MustCompile(r.pattern).ReplaceAllString(text, r.replace)
	}

	text = collapseRepeatedPunctuation(text)

	// Remove boilerplate
	boilerplates := []string{
		"accept cookies", "cookie policy", "privacy policy", "all rights reserved",
		"terms and conditions", "share this", "follow us", "sign up for",
	}
	for _, bp := range boilerplates {
		text = strings.ReplaceAll(text, bp, "")
	}

	if len(text) > 1000 {
		text = text[:997] + "..."
	}

	return strings.TrimSpace(text)
}
