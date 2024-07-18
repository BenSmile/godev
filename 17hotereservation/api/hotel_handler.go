package api

import (
	"github.com/bensmile/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	roomStore  db.RoomStore
	hotelStore db.HotelStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		roomStore:  rs,
		hotelStore: hs,
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

	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)

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

	rooms, err := h.roomStore.GetRooms(c.Context(), filter)

	if err != nil {
		return err
	}

	return c.JSON(rooms)

}
