package presenter

import (
	"code-arcade-web-search-api/pkg/search"

	"github.com/gofiber/fiber/v2"
)

type SearchMeta struct {
	TimeTaken string `json:"timeTaken"`
	Length    int    `json:"length"`
	Size      string `json:"size"`
}

func WebSearchSuccessResponse(results []search.ScrapedPage, meta SearchMeta) *fiber.Map {
	return &fiber.Map{
		"success": true,
		"meta":    meta,
		"results": results,
	}
}

func WebSearchErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"success": false,
		"meta":    nil,
		"results": nil,
		"error":   err.Error(),
	}
}
