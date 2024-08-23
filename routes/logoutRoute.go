package routes

import (
	"net/http"

	"github.com/froggy-12/mooshroombase/internal/types"
	"github.com/gofiber/fiber/v2"
)

func LogOut(c *fiber.Ctx) error {
	cookie := &fiber.Cookie{
		Name:     "jwtToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(cookie)
	return c.Status(http.StatusAccepted).JSON(types.HttpSuccessResponse{Message: "user has been logged out successfully"})
}
