package routes

import (
	"github.com/froggy-12/mooshroombase/internal/types"
	"github.com/gofiber/fiber/v2"
)

func FeaturedRoutes(router fiber.Router) {
	router.Get("/ping", pong)
}

func pong(c *fiber.Ctx) error {
	return c.JSON(types.HttpSuccessResponse{Message: "pong"})
}
