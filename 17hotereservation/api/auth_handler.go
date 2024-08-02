package api

import (
	"fmt"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {

	var authParams AuthParams

	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authParams.Password))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}

	fmt.Printf("authenticated -> : %+v\n", user)

	return nil
}
