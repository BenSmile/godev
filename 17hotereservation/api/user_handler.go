package api

import (
	"errors"
	"fmt"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlerGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleLogin(c *fiber.Ctx) error {
	type AuthParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var authParams AuthParams

	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	fmt.Println("body : ", authParams)

	return nil
}

func (h *UserHandler) HandlerDeleteuser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = c.Context()
	)
	err := h.userStore.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{
		"deleted": id,
	})
}

func (h *UserHandler) HandlerCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) != 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	newUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(newUser)
}

func (h *UserHandler) HandlerUpdateUser(c *fiber.Ctx) error {
	var (
		id     = c.Params("id")
		params types.UpdateUserParams
		ctx    = c.Context()
	)
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	_, err := h.userStore.UpdateUser(ctx, id, &params)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{
		"updated": id,
	})
}

func (h *UserHandler) HandlerGetUserByID(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = c.Context()
	)
	user, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{
				"message": "item not found for the given id",
			})
		}
		return err
	}
	return c.JSON(user)
}
