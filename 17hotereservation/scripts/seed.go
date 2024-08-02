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
	userStore  db.MongoUserStore
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
	userStore = *db.NewMongoUserStore(client)
}

func seedUser(fname, lname, email string) {

	user, err := types.NewUserFromParams(
		types.CreateUserParams{
			Password:  "pass",
			FirstName: fname,
			LastName:  lname,
			Email:     email},
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = userStore.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

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
	seedUser("admin", "admin", "admin@test.com")
}
