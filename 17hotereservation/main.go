package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/bensmile/hotel-reservation/api"
	"github.com/bensmile/hotel-reservation/db"
	"github.com/bensmile/hotel-reservation/middleware"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
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
		userStore    = db.NewMongoUserStore(client)
		userHandler  = api.NewUserHandler(userStore)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		store        = db.Store{
			Hotel:   hotelStore,
			User:    userStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
		hotelHandler   = api.NewHotelHandler(&store)
		bookingHandler = api.NewBookingHandler(&store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(&store)
		app            = fiber.New(config)
		apiV1          = app.Group("/api/v1", middleware.JWTAuth(userStore))
		// admin          = apiV1.Group("/admin", middleware.OnlyAdmin)
		admin = apiV1.Group("/admin")
	)

	// users
	apiV1.Get("users", userHandler.HandlerGetUsers)
	apiV1.Get("users/:id", userHandler.HandlerGetUserByID)
	apiV1.Post("users", userHandler.HandlerCreateUser)
	apiV1.Put("users/:id", userHandler.HandlerUpdateUser)
	apiV1.Delete("users/:id", userHandler.HandlerDeleteuser)

	// hotels
	apiV1.Get("hotels", hotelHandler.HandleGetHotels)
	apiV1.Get("hotels/:id/rooms", hotelHandler.HandleGetRoomByHotel)
	apiV1.Get("hotels/:id", hotelHandler.HandleGetHotelById)

	// rooms
	apiV1.Get("rooms", roomHandler.HandleGetRooms)
	apiV1.Post("rooms/:id/book", roomHandler.HandleBookRoom)

	// booking
	admin.Get("bookings", bookingHandler.HandleGetBookings)
	apiV1.Put("bookings/:id/cancel", bookingHandler.HandleCancelBooking)
	admin.Get("bookings/:id", bookingHandler.HandleGetBookingById)
	apiV1.Get("bookings", bookingHandler.HandleGetBookingsForUser)

	app.Get("/", handleHome)
	app.Post("/api/auth", authHandler.HandleLogin)

	app.Listen(*listenAddr)

}

func handleHome(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Running....",
	})
}
