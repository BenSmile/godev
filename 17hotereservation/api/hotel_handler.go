package api

import (
	"fmt"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
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

	type HotelQueryParams struct {
		Rooms  bool
		Rating int
	}

	var qParams HotelQueryParams

	if err := c.QueryParser(&qParams); err != nil {
		return err
	}

	fmt.Println(qParams)

	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)

	if err != nil {
		return err
	}

	return c.JSON(hotels)
}
