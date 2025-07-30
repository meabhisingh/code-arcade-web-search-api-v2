package handlers

import (
	"code-arcade-web-search-api/api/presenter"
	"code-arcade-web-search-api/config"
	"code-arcade-web-search-api/pkg/search"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SearchRequest struct {
	Query string `json:"query"`
}

func WebSearch(service search.Service) fiber.Handler {

	return func(c *fiber.Ctx) error {

		var body SearchRequest

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid JSON body",
			})
		}

		if body.Query == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Please enter a query",
			})
		}

		start := time.Now()

		searchResult, err := service.FetchSearchResults(body.Query)
		if len(searchResult) > 10 {
			searchResult = searchResult[:10]
		}

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": fmt.Sprintf("Error while fetching: %v", err),
			})
		}

		results, err := service.Scrape(searchResult)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": fmt.Sprintf("Error while fetching: %v", err),
			})
		}

		totalContent := ""

		for _, curr := range results {
			totalContent += curr.Content
		}

		byteSize := len(totalContent)
		kbSizeFloat := float64(byteSize) / 1024.0
		kbSizeFormatted := fmt.Sprintf("%.2f", kbSizeFloat)

		end := time.Now()

		duration := end.Sub(start)
		durationInSeconds := int(duration.Seconds())

		meta := presenter.SearchMeta{
			TimeTaken: fmt.Sprintf("%ds", durationInSeconds),
			Length:    len(totalContent),
			Size:      fmt.Sprintf("%skb", kbSizeFormatted),
		}

		return c.JSON(presenter.WebSearchSuccessResponse(results, meta))
	}
}

func TestSearchEngine() fiber.Handler {

	return func(c *fiber.Ctx) error {

		searxngUrl := config.GetEnv("SEARCH_ENGINE_URL", "http://localhost:8080")

		resp, err := http.Get(searxngUrl)

		if err != nil || resp.StatusCode != http.StatusOK {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Search Engine not working",
			})
		}

		defer resp.Body.Close()

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Test route works!",
		})
	}
}
