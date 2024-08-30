package middleware

import (
	"fmt"

	"github.com/bensmile/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func OnlyAdmin(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return fmt.Errorf("not authorized")
	}

	if !user.IsAdmin {
		return fmt.Errorf("not authorized")
	}
	return c.Next()
}
