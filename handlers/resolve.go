package handlers

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/peyzoman/url-shortner/database"
)

func Resolve(ctx *fiber.Ctx) error {
	url := ctx.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	} else if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	r1 := database.CreateClient(1)
	defer r1.Close()

	r1.Incr(database.Ctx, "counter")
	return ctx.Redirect(value, http.StatusMovedPermanently)
}
