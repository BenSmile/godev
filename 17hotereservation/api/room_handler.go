package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// DATE_FORMAT = time.RFC3339
	DATE_FORMAT = "2006-01-02 15:04:05"
)

type BookingRoomReq struct {
	FromDate        string `json:"fromDate"`
	TillDate        string `json:"tillDate"`
	ContactNumber   string `json:"contactNumber"`
	NumberOfPersons int    `json:"numberOfPersons"`
}

func validateBookingDates(fromDate, tillDate string) (time.Time, time.Time, error) {

	fmt.Println("date : ", fromDate)

	from, err := time.Parse(DATE_FORMAT, fromDate+" 00:00:00")
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid fromDate format, expected yyyy-MM-dd")
	}

	till, err := time.Parse(DATE_FORMAT, tillDate+" 00:00:00")
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid tillDate format, expected yyyy-MM-dd")
	}

	if from.After(till) {
		return time.Time{}, time.Time{}, fmt.Errorf("fromDate must be before tillDate")
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

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	roomOID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return fmt.Errorf("invalid room id")
	}

	var bookingRoomBody BookingRoomReq

	if err := c.BodyParser(&bookingRoomBody); err != nil {
		return err
	}

	fmt.Printf("bookingRoomBody : %+v\n", bookingRoomBody)

	from, to, err := validateBookingDates(bookingRoomBody.FromDate, bookingRoomBody.TillDate)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	userID, ok := c.Context().UserValue("user").(string)
	if !ok {
		return fmt.Errorf("invalid user id")
	}
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}

	booking := types.Booking{
		RoomID:          roomOID,
		UserID:          userOID,
		FromDate:        from,
		TillDate:        to,
		NumberOfPersons: bookingRoomBody.NumberOfPersons,
	}

	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(insertedBooking)
}
