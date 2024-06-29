package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))

	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "France",
	}

	insertedHotel, err := hotelStore.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedHotel)

	rooms := []types.Room{
		{
			Types:     types.SingleRoomType,
			BasePrice: 88.9,
			HotelID:   hotel.ID,
		}, {
			Types:     types.SingleRoomType,
			BasePrice: 88.9,
			HotelID:   hotel.ID,
		},
	}

	for _, room := range rooms {
		insertedRoom, err := roomStore.InsertRoom(context.TODO(), &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}

}
