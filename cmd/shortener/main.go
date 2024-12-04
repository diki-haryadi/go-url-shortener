package main

import (
	shortener "go-url-shortener"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", shortener.Resolve)
	app.Post("/api/v1", shortener.Shorten)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment file.")
	}

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	app.Listen(":3000") // + os.Getenv("APP_PORT"))
	// base62EncodedString := helpers.Base62Encode(9999999)
	// fmt.Println(base62EncodedString)
	// fmt.Println(helpers.Base62Decode(base62EncodedString))
}
