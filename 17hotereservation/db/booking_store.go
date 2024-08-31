package db

import (
	"context"
	"fmt"

	"github.com/bensmile/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	bookingCollection = "bookings"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]types.Booking, error)
	GetBookingById(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M) (*types.Booking, error)
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]types.Booking, error) {
	curs, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []types.Booking

	if err := curs.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {

	res, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(bookingCollection),
	}
}

type MongoBookingStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, values bson.M) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %w", err)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: values}}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("booking with ID %s not found", id)
	}

	var updatedBooking types.Booking
	err = s.collection.FindOne(ctx, filter).Decode(&updatedBooking)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated booking: %w", err)
	}

	return &updatedBooking, nil
}
