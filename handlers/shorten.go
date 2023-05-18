package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL      string        `json:"url"`
	ShortURL string        `json:"short_url"`
	Expiry   time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	ShortURL        string        `json:"short_url"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenHandler(ctx *fiber.Ctx) error {
	body := &request{}

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse request body"})
	}

	resp := &response{
		URL:             body.URL,
		ShortURL:        body.ShortURL,
		Expiry:          body.Expiry,
		XRateRemaining:  5,
		XRateLimitReset: 30,
	}
	return ctx.Status(http.StatusOK).JSON(resp)
}
