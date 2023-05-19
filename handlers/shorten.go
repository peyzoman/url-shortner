package handlers

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/peyzoman/url-shortner/database"
	"github.com/peyzoman/url-shortner/utils"
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

	// rate limit implementation
	r1 := database.CreateClient(1)
	defer r1.Close()

	val, err := r1.Get(database.Ctx, ctx.IP()).Result()
	limit, _ := r1.TTL(database.Ctx, ctx.IP()).Result()

	if err == redis.Nil {
		_ = r1.Set(database.Ctx, ctx.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else if err != nil {
		return err
	}

	valInt, _ := strconv.Atoi(val)
	if valInt <= 0 {
		return ctx.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error":            "rate limit exceeded",
			"rate_limit_reset": limit / time.Nanosecond / time.Minute,
		})
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(body.URL) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	if !utils.RemoveDomainError(body.URL) {
		return ctx.Status(http.StatusServiceUnavailable).JSON(fiber.Map{"error": "can't do that"})
	}

	var id string
	if body.ShortURL == "" {
		id = utils.Base62Encode(rand.Uint64())
	} else {
		id = body.ShortURL
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "URL short is already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to connect",
		})
	}

	defaultAPIQuotaStr := os.Getenv("API_QUOTA")
	defaultAPIQuota, _ := strconv.Atoi(defaultAPIQuotaStr)

	resp := &response{
		URL:             body.URL,
		ShortURL:        "",
		Expiry:          body.Expiry,
		XRateRemaining:  defaultAPIQuota,
		XRateLimitReset: 30,
	}

	remainingQuota, err := r1.Decr(database.Ctx, ctx.IP()).Result()
	if err != nil {
		return err
	}

	resp.XRateRemaining = int(remainingQuota)
	resp.XRateRemaining = int(limit / time.Nanosecond / time.Minute)
	resp.ShortURL = os.Getenv("DOMAIN") + "/" + id

	return ctx.Status(http.StatusOK).JSON(resp)
}
