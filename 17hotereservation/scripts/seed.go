package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.MongoRoomStore
	hotelStore db.MongoHotelStore
	ctx        = context.Background()
)

func init() {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(db.DB_URI))

	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = *db.NewMongoHotelStore(client)
	roomStore = *db.NewMongoRoomStore(client, &hotelStore)
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedHotel)

	rooms := []types.Room{
		{
			Size:      "small",
			BasePrice: 80.9,
			HotelID:   hotel.ID,
		}, {
			Size:      "normal",
			BasePrice: 100.9,
			HotelID:   hotel.ID,
		}, {
			Size:      "kingsize",
			BasePrice: 120.9,
			HotelID:   hotel.ID,
		},
	}

	for _, room := range rooms {
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}

}

func main() {
	seedHotel("Serena", "Congo", 5)
	seedHotel("Panorama", "Congo", 4)
}
