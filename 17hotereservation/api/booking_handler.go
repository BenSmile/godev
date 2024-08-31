package api

import (
	"fmt"
	"net/http"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/helpers"
	"github.com/bensmile/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// only for admin
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {

	var bookings []types.Booking

	connectedUser, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	filter := bson.M{}

	if !connectedUser.IsAdmin {
		filter = bson.M{"userID": connectedUser.ID}
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(bookings)
}

// only for admin
func (h *BookingHandler) HandleGetBookingById(c *fiber.Ctx) error {

	booking, err := h.store.Booking.GetBookingById(c.Context(), c.Params("id"))

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {

	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return err
	}

	user, err := helpers.GetAuthUser(c)
	fmt.Printf("User : %+v\n", user)
	fmt.Printf("Booking : %+v\n", booking)
	fmt.Printf("Equals : %+v\n", user.ID != booking.UserID || !user.IsAdmin)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	if user.ID != booking.UserID && !user.IsAdmin {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	if booking.Canceled {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "already canceled"})
	}

	booking, err = h.store.Booking.UpdateBooking(c.Context(), id, bson.M{"canceld": true})

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(booking)
}

// only for connected user
func (h *BookingHandler) HandleGetBookingsForUser(c *fiber.Ctx) error {

	user, err := helpers.GetAuthUser(c)

	if err != nil {
		return err
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{
		"userID": user.ID,
	})

	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(bookings)
}
