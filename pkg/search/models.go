package search

type SearchResult struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	PublishedDate string `json:"publishedDate"`
}

type ScrapedPageSource struct {
	Url     string `json:"url"`
	Title   string `json:"title"`
	Website string `json:"website"`
	Icon    string `json:"icon"`
}

type ScrapedPage struct {
	Source  ScrapedPageSource `json:"source"`
	Content string            `json:"content"`
}

type SearchEngineResponse struct {
	Query               string         `json:"query"`
	NumberOfResults     int            `json:"number_of_results"`
	Results             []SearchResult `json:"results"`
	Answers             []Answer       `json:"answers"`
	Corrections         []interface{}  `json:"corrections"`
	Infoboxes           []interface{}  `json:"infoboxes"`
	Suggestions         []string       `json:"suggestions"`
	UnresponsiveEngines [][]string     `json:"unresponsive_engines"`
}

type Answer struct {
	URL       string        `json:"url"`
	Template  string        `json:"template"`
	Engine    string        `json:"engine"`
	ParsedURL []interface{} `json:"parsed_url"`
	Answer    string        `json:"answer"`
}
