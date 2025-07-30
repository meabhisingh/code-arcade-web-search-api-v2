package main

import (
	"code-arcade-web-search-api/api/routes"
	"code-arcade-web-search-api/config"
	"code-arcade-web-search-api/pkg/search"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	config.LoadEnv()
	app := fiber.New()

	port := config.GetEnv("PORT", "3000")

	app.Static("/", "./public")

	log.Printf("Server starting on port %s...", port)

	searchService := search.NewService()
	// Registering Routes
	routes.SearchRouter(app.Group("/api/v1"), searchService)

	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	// Listen on that port
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}

}
