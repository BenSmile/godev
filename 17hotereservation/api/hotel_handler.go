package api

import (
	"github.com/bensmile/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	// type HotelQueryParams struct {
	// 	Rooms  bool
	// 	Rating int
	// }

	// var qParams HotelQueryParams

	// if err := c.QueryParser(&qParams); err != nil {
	// 	return err
	// }

	// fmt.Printf("%+v", qParams)

	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)

	if err != nil {
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetRoomByHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"hotel_id": oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filter)

	if err != nil {
		return err
	}

	return c.JSON(rooms)

}

func (h *HotelHandler) HandleGetHotelById(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	hotel, err := h.store.Hotel.GetHotelById(c.Context(), oid)

	if err != nil {
		return err
		// return c.JSON(map[string]string{
		// 	"message": err.Error(),
		// })
	}

	return c.JSON(hotel)

}
