package api

import (
	"net/http"

	"github.com/bensmile/hotel-reservation/db"
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

	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})

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

	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})

	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(bookings)
}

// only for connected user
func (h *BookingHandler) HandleGetBookingsForUser(c *fiber.Ctx) error {

	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{
		"userID": user.ID,
	})

	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(bookings)
}
