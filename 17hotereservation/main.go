package main

import (
	"context"
	"flag"
	"log"

	"github.com/bensmile/hotel-reservation/api"
	"github.com/bensmile/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{
			"error": err.Error(),
		})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":3100", "The listen address of the api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))

	if err != nil {
		log.Fatal(err)
	}
	// handler initialization
	var (
		userHandler  = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
		app          = fiber.New(config)
		apiV1        = app.Group("/api/v1")
	)

	apiV1.Get("users", userHandler.HandlerGetUsers)
	apiV1.Get("users/:id", userHandler.HandlerGetUserByID)
	apiV1.Post("users", userHandler.HandlerCreateUser)
	apiV1.Put("users/:id", userHandler.HandlerUpdateUser)
	apiV1.Delete("users/:id", userHandler.HandlerDeleteuser)
	apiV1.Get("hotels", hotelHandler.HandleGetHotels)
	app.Get("/", handleHome)
	app.Listen(*listenAddr)

}

func handleHome(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"message": "Running...",
	})
}
