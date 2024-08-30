package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int                  `bsom:"rating" json:"rating"`
}

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BasePrice float64            `bson:"base_price" json:"basePrice"`
	// small, normal, kingsize
	Size     string             `bson:"size" json:"size"`
	SeaSide  bool               `bson:"sea_side" json:"seaSide"`
	IsBooked bool               `bson:"-" json:"isBooked"`
	Price    float64            `bson:"price" json:"price"`
	HotelID  primitive.ObjectID `bson:"hotel_id" json:"hotelID"`
}

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeaSideleRoomType
	DeluxeRoomType
)

type UpdateRoomsParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
