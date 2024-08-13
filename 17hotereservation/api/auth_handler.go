package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {

	var authParams types.AuthParams

	if err := c.BodyParser(&authParams); err != nil {
		// return err
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid credentials",
		})
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		// return fmt.Errorf("invalid credentials")
		// return c.Status(http.StatusUnauthorized).SendString("invalid credentials")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid credentials",
		})
	}

	if !types.IsValidPassword(user.Password, authParams.Password) {
		// return fmt.Errorf("invalid credentials")
		// return c.Status(http.StatusUnauthorized).SendString("invalid credentials")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid credentials",
		})

	}

	return c.JSON(types.AuthResponse{
		User:  user,
		Token: makeClaimsFromuser(user),
	})

}

func makeClaimsFromuser(user *types.User) string {
	now := time.Now()
	expiredAt := now.Add(time.Hour * 24)
	claims := jwt.MapClaims{
		"userID":    user.ID,
		"email":     user.Email,
		"expiredAt": expiredAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret : ", err)
	}
	return tokenString
}
