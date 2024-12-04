package shortener

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"go-url-shortener/pkg/cache"
)

func Resolve(ctx *fiber.Ctx) error {
	url := ctx.Params("url")

	r := cache.CreateClient(0)
	defer r.Close()

	value, err := r.Get(cache.Ctx, url).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short-url not found in db"})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal error"})
	}

	rInr := cache.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(cache.Ctx, "counter")

	return ctx.Redirect(value, 301)
}
