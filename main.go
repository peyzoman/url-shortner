package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/peyzoman/url-shortner/handlers"
)

func setupRoutes(app *fiber.App) {
	app.Post("/api/v1/shorten", handlers.ShortenHandler)
	app.Get("/:url", handlers.Resolve)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load environment file")
	}

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	app.Listen(os.Getenv("APP_PORT"))
}
