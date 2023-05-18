package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load environment file")
	}

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	app.Listen(":3000")
}
