package search

type Service interface {
	FetchSearchResults(query string) ([]SearchResult, error)
	Scrape(results []SearchResult) ([]ScrapedPage, error)
}

type service struct {
	// add clients here if needed, e.g. httpClient
}

func NewService() Service {
	return &service{}
}
