package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DATE_FORMAT = "2006-01-02 15:04:05"
)

type BookingRoomReq struct {
	FromDate        string `json:"fromDate"`
	TillDate        string `json:"tillDate"`
	ContactNumber   string `json:"contactNumber"`
	NumberOfPersons int    `json:"numberOfPersons"`
}

func validateBookingDates(fromDate, tillDate string) (time.Time, time.Time, error) {

	var nilDate time.Time

	currentTime := time.Now().Truncate(24 * time.Hour) // Truncate to remove the time part

	from, err := time.Parse(DATE_FORMAT, fmt.Sprintf("%s 00:00:00", fromDate))
	if err != nil {
		return nilDate, nilDate, fmt.Errorf("invalid fromDate format, expected yyyy-MM-dd")
	}

	till, err := time.Parse(DATE_FORMAT, fmt.Sprintf("%s 00:00:00", tillDate))
	if err != nil {
		return nilDate, nilDate, fmt.Errorf("invalid tillDate format, expected yyyy-MM-dd")
	}

	if from.After(till) {
		return nilDate, nilDate, fmt.Errorf("fromDate must be before tillDate")
	}

	if from.Before(currentTime) {
		return nilDate, nilDate, fmt.Errorf("fromDate cannot be before today")
	}

	if till.Before(currentTime) {
		return nilDate, nilDate, fmt.Errorf("tillDate cannot be before today")
	}

	return from, till, nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {

	var bookingRoomBody BookingRoomReq

	if err := c.BodyParser(&bookingRoomBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	from, to, err := validateBookingDates(bookingRoomBody.FromDate, bookingRoomBody.TillDate)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	roomIdStr := c.Params("id")

	room, err := h.store.Room.GetRoomById(c.Context(), roomIdStr)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("room not found for the given id : %s", roomIdStr)})
	}

	roomAvailable, _ := h.isRoomAvailable(c.Context(), room.ID, from, to)

	if !roomAvailable {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "room already booked"})
	}

	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "unauthorized"})
	}

	booking := types.Booking{
		RoomID:          room.ID,
		UserID:          user.ID,
		FromDate:        from,
		TillDate:        to,
		NumberOfPersons: bookingRoomBody.NumberOfPersons,
	}

	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	// room.IsBooked = true

	// values := bson.M{
	// 	"is_booked": true,
	// }

	// h.store.Room.UpdateRoom(c.Context(), roomIdStr, values)

	return c.Status(http.StatusCreated).JSON(insertedBooking)
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, oid primitive.ObjectID, from time.Time, to time.Time) (bool, error) {
	filter := bson.M{
		"fromDate": bson.M{
			"$gte": from,
		},
		"tillDate": bson.M{
			"$lte": to,
		},
		"roomID": oid,
	}
	bookings, err := h.store.Booking.GetBookings(ctx, filter)

	if err != nil {
		return false, err
	}

	return len(bookings) == 0, nil
}
