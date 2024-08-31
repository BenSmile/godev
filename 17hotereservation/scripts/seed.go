package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bensmile/hotel-reservation/api"
	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.MongoRoomStore
	bookingStore db.MongoBookingStore
	userStore    db.MongoUserStore
	hotelStore   db.MongoHotelStore
	ctx          = context.Background()
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
	bookingStore = *db.NewMongoBookingStore(client)
}

func seedUser(fname, lname, email string, isAdmin bool) *types.User {

	user, err := types.NewUserFromParams(
		types.CreateUserParams{
			Password:  "pass",
			FirstName: fname,
			LastName:  lname,
			Email:     email,
			IsAdmin:   isAdmin},
	)
	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := userStore.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s\n", email, api.MakeClaimsFromuser(user))

	return insertedUser

}

func seedHotel(name, location string, rating int) *types.Hotel {
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
		fmt.Printf("Room : %+v\n", insertedRoom)
	}

	return insertedHotel
}

func seedRoom(ss string, bp float64, hotel primitive.ObjectID) *types.Room {

	room := types.Room{
		Size:      ss,
		BasePrice: bp,
		HotelID:   hotel,
	}

	insertedRoom, err := roomStore.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Room : %+v\n", insertedRoom)

	return insertedRoom

}

func seedBookingRoom(userID, roomID primitive.ObjectID, numberOfPersons int, from, to time.Time, canceled bool) *types.Booking {

	booking := types.Booking{
		RoomID:          roomID,
		UserID:          userID,
		NumberOfPersons: numberOfPersons,
		FromDate:        from,
		TillDate:        to,
	}

	insertedBooking, err := bookingStore.InsertBooking(ctx, &booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Room : %+v\n", insertedBooking)

	return insertedBooking

}

func main() {
	hotel1 := seedHotel("Serena", "Congo", 5)

	room1 := seedRoom("small", 99.9, hotel1.ID)

	admin := seedUser("admin", "admin", "admin@test.com", true)

	user := seedUser("admin2", "admin", "admin2@test.com", false)

	booking := seedBookingRoom(user.ID, room1.ID, 3, time.Now(), time.Now().AddDate(0, 0, 2), false)

	booking2 := seedBookingRoom(admin.ID, room1.ID, 1, time.Now().AddDate(0, 0, 4), time.Now().AddDate(0, 0, 5), false)

	fmt.Printf("[]booking : %+v\n", booking)

	fmt.Printf("[]booking2 : %+v\n", booking2)

}
