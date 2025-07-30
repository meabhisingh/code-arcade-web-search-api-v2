package routes

import (
	"code-arcade-web-search-api/api/handlers"
	"code-arcade-web-search-api/pkg/search"

	"github.com/gofiber/fiber/v2"
)

func SearchRouter(app fiber.Router, service search.Service) {

	app.Post("/search", handlers.WebSearch(service))
	app.Get("/test", handlers.TestSearchEngine())

}
